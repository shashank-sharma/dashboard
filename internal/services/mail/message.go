package mail

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/pocketbase/pocketbase/tools/types"
	"google.golang.org/api/gmail/v1"

	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
)

func (ms *MailService) SyncMessages(mailSync *models.MailSync) error {
	client, err := ms.FetchClient(mailSync.Token)
	if err != nil {
		return err
	}

	gmailService, err := ms.GetGmailService(client)
	if err != nil {
		logger.LogError("Unable to create Gmail service")
		return err
	}

	if mailSync.LastSyncState != "" {
		if err := ms.incrementalSync(gmailService, mailSync); err == nil {
			return nil
		}
		logger.LogInfo("Falling back to full sync")
	}

	return ms.fullSync(gmailService, mailSync)
}

func (ms *MailService) fullSync(srv *gmail.Service, mailSync *models.MailSync) error {
	messagesChannel := make(chan *gmail.Message, 100)
	errorChannel := make(chan error, 1)
	var wg sync.WaitGroup

	const numWorkers = 5
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range messagesChannel {
				if err := ms.processMessage(srv, msg.Id, mailSync); err != nil {
					select {
					case errorChannel <- err:
					default:
					}
				}
			}
		}()
	}

	go func() {
		defer close(messagesChannel)
		pageToken := ""

		for {
			// TODO: Testing with sample labels
			req := srv.Users.Messages.List("me").LabelIds("Label_6").MaxResults(100)
			if pageToken != "" {
				req.PageToken(pageToken)
			}

			resp, err := req.Do()
			if err != nil {
				select {
				case errorChannel <- err:
				default:
				}
				return
			}

			for _, msg := range resp.Messages {
				messagesChannel <- msg
			}

			pageToken = resp.NextPageToken
			if pageToken == "" {
				break
			}
		}
	}()

	wg.Wait()
	close(errorChannel)

	if err := <-errorChannel; err != nil {
		return fmt.Errorf("error during sync: %v", err)
	}

	profile, err := srv.Users.GetProfile("me").Do()
	if err != nil {
		return fmt.Errorf("failed to get profile: %v", err)
	}

	return query.UpdateRecord[*models.MailSync](mailSync.Id, map[string]interface{}{
		"last_sync_state": strconv.FormatUint(profile.HistoryId, 10),
		"last_synced":     types.NowDateTime(),
	})
}

func (ms *MailService) processMessage(srv *gmail.Service, messageId string, mailSync *models.MailSync) error {
	msg, err := srv.Users.Messages.Get("me", messageId).Format("full").Do()
	if err != nil {
		return fmt.Errorf("failed to get message %s: %v", messageId, err)
	}

	body := ms.extractMessageBody(msg.Payload)
	from, to := ms.extractEmailAddresses(msg.Payload.Headers)
	subject := ms.extractHeader(msg.Payload.Headers, "Subject")

	internalDate := types.DateTime{}
	receivedDate := types.DateTime{}
	internalDate.Scan(time.Unix(0, msg.InternalDate*int64(time.Millisecond)))

	if receivedStr := ms.extractHeader(msg.Payload.Headers, "Received"); receivedStr != "" {
		if parsedTime, err := time.Parse(time.RFC1123Z, receivedStr); err == nil {
			receivedDate.Scan(parsedTime)
		}
	}
	if receivedDate.IsZero() {
		receivedDate = internalDate
	}

	isUnread := false
	isImportant := false
	isStarred := false
	isSpam := false
	isInbox := false
	isTrash := false
	isDraft := false
	isSent := false

	for _, labelId := range msg.LabelIds {
		switch labelId {
		case "UNREAD":
			isUnread = true
		case "IMPORTANT":
			isImportant = true
		case "STARRED":
			isStarred = true
		case "SPAM":
			isSpam = true
		case "INBOX":
			isInbox = true
		case "TRASH":
			isTrash = true
		case "DRAFT":
			isDraft = true
		case "SENT":
			isSent = true
		}
	}

	externalData := map[string]interface{}{
		"history_id":    msg.HistoryId,
		"label_ids":     msg.LabelIds,
		"size_estimate": msg.SizeEstimate,
	}

	externalDataStr, err := json.Marshal(externalData)
	if err != nil {
		// Handle error
		return err
	}

	mailMessage := &models.MailMessage{
		User:         mailSync.User,
		MailSync:     mailSync.Id,
		MessageId:    msg.Id,
		ThreadId:     msg.ThreadId,
		From:         from,
		To:           to,
		Subject:      subject,
		Snippet:      msg.Snippet,
		Body:         body,
		IsUnread:     isUnread,
		IsImportant:  isImportant,
		IsStarred:    isStarred,
		IsSpam:       isSpam,
		IsInbox:      isInbox,
		IsTrash:      isTrash,
		IsDraft:      isDraft,
		IsSent:       isSent,
		InternalDate: internalDate,
		ReceivedDate: receivedDate,
		ExternalData: string(externalDataStr),
	}

	return query.UpsertRecord[*models.MailMessage](mailMessage, map[string]interface{}{
		"message_id": msg.Id,
	})
}

func (ms *MailService) incrementalSync(srv *gmail.Service, mailSync *models.MailSync) error {
	historyId, err := strconv.ParseUint(mailSync.LastSyncState, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid history ID: %v", err)
	}

	req := srv.Users.History.List("me").StartHistoryId(historyId)
	var processedMsgIds = make(map[string]bool)
	pageToken := ""

	for {
		req.PageToken(pageToken)
		resp, err := req.Do()
		if err != nil {
			return fmt.Errorf("failed to get history: %v", err)
		}

		for _, history := range resp.History {
			for _, added := range history.MessagesAdded {
				if !processedMsgIds[added.Message.Id] {
					if err := ms.processMessage(srv, added.Message.Id, mailSync); err != nil {
						logger.LogError("Error processing message", err)
					}
					processedMsgIds[added.Message.Id] = true
				}
			}

			for _, labelChanged := range history.LabelsAdded {
				if !processedMsgIds[labelChanged.Message.Id] {
					if err := ms.processMessage(srv, labelChanged.Message.Id, mailSync); err != nil {
						logger.LogError("Error processing label change", err)
					}
					processedMsgIds[labelChanged.Message.Id] = true
				}
			}
		}

		pageToken = resp.NextPageToken
		if pageToken == "" {
			break
		}
	}

	profile, err := srv.Users.GetProfile("me").Do()
	if err != nil {
		return fmt.Errorf("failed to get profile: %v", err)
	}

	return query.UpdateRecord[*models.MailSync](mailSync.Id, map[string]interface{}{
		"last_sync_state": strconv.FormatUint(profile.HistoryId, 10),
		"last_synced":     types.NowDateTime(),
	})
}

func (ms *MailService) extractMessageBody(payload *gmail.MessagePart) string {
	if payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err != nil {
			return ""
		}
		return string(data)
	}

	if len(payload.Parts) > 0 {
		var htmlPart, plainPart string
		for _, part := range payload.Parts {
			switch part.MimeType {
			case "text/html":
				htmlPart = ms.extractMessageBody(part)
			case "text/plain":
				plainPart = ms.extractMessageBody(part)
			}
		}
		if htmlPart != "" {
			return htmlPart
		}
		return plainPart
	}

	return ""
}

func (ms *MailService) extractEmailAddresses(headers []*gmail.MessagePartHeader) (string, string) {
	var from, to string
	for _, header := range headers {
		switch header.Name {
		case "From":
			from = header.Value
		case "To":
			to = header.Value
		}
	}
	return from, to
}

func (ms *MailService) extractHeader(headers []*gmail.MessagePartHeader, name string) string {
	for _, header := range headers {
		if header.Name == name {
			return header.Value
		}
	}
	return ""
}
