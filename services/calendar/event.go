package calendar

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/shashank-sharma/backend/logger"
	"github.com/shashank-sharma/backend/models"
	"github.com/shashank-sharma/backend/query"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

func GetClient(calTokenId string) (*http.Client, error) {
	googleConfig := LoadCalendarConfig()
	calendarToken, err := query.FindById[*models.CalendarToken](calTokenId)
	if err != nil {
		return nil, err
	}

	oauthToken := &oauth2.Token{
		AccessToken:  calendarToken.AccessToken,
		TokenType:    calendarToken.TokenType,
		RefreshToken: calendarToken.RefreshToken,
		Expiry:       calendarToken.Expiry.Time(),
	}

	return googleConfig.Client(context.Background(), oauthToken), nil
}

func SyncEvents(calendarSync *models.CalendarSync) error {
	client, err := GetClient(calendarSync.Token)
	if err != nil {
		return err
	}
	calendarService, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))

	if err != nil {
		logger.LogError("Unable to create calendar service")
		return err
	}

	request := calendarService.Events.List(calendarSync.Type)
	if calendarSync.SyncToken != "" {
		request.SyncToken(calendarSync.SyncToken)
	} else {
		timeMin := time.Now().AddDate(0, -2, 0).Format(time.RFC3339)
		timeMax := time.Now().AddDate(0, 2, 0).Format(time.RFC3339)

		request.TimeMin(timeMin)
		request.TimeMax(timeMax)
	}

	eventsChannel := make(chan *calendar.Event)
	var wg sync.WaitGroup

	const numWorkers = 5
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for event := range eventsChannel {
				InsertEvent(event, calendarSync.User, calendarSync.Id)
			}
		}()
	}

	pageToken := ""
	for {
		request.PageToken(pageToken)
		events, err := request.SingleEvents(true).Do()
		if err != nil {
			if e, ok := err.(*googleapi.Error); ok && e.Code == 410 {
				logger.LogError("Invalid sync token")
				return err
			}
			return err
		}
		if len(events.Items) == 0 {
			logger.Debug.Println("No new events to sync")
			break
		}

		for _, event := range events.Items {
			eventsChannel <- event
		}

		pageToken = events.NextPageToken
		if pageToken == "" {
			query.UpdateRecord[*models.CalendarSync](calendarSync.Id, map[string]interface{}{
				"sync_token":  events.NextSyncToken,
				"last_synced": types.NowDateTime(),
			})
			break
		}
	}

	close(eventsChannel)
	wg.Wait()

	return nil
}

func InsertEvent(event *calendar.Event, userId string, calendarSyncId string) error {

	eventModel := &models.CalendarEvent{
		CalendarId:     event.Id,
		CalendarUId:    event.ICalUID,
		User:           userId,
		Calendar:       calendarSyncId,
		Etag:           event.Etag,
		Summary:        event.Summary,
		Description:    event.Description,
		EventType:      event.EventType,
		Creator:        event.Creator.DisplayName,
		CreatorEmail:   event.Creator.Email,
		Organizer:      event.Organizer.DisplayName,
		OrganizerEmail: event.Organizer.Email,
		Kind:           event.Kind,
		Location:       event.Location,
		Status:         event.Status,
	}

	if calendarStart, err := types.ParseDateTime(event.Start.DateTime); err == nil {
		eventModel.Start = calendarStart
	}
	if calendarEnd, err := types.ParseDateTime(event.End.DateTime); err == nil {
		eventModel.End = calendarEnd
	}
	if calendarEventCreated, err := types.ParseDateTime(event.Created); err == nil {
		eventModel.EventCreated = calendarEventCreated
	}
	if calendarEventUpdated, err := types.ParseDateTime(event.Updated); err == nil {
		eventModel.EventCreated = calendarEventUpdated
	}

	query.UpsertRecord[*models.CalendarEvent](eventModel, map[string]interface{}{
		"calendar_id": event.Id,
	})

	return nil
}
