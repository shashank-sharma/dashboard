package routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
	"github.com/shashank-sharma/backend/internal/services/mail"
	"github.com/shashank-sharma/backend/internal/services/oauth"
	"github.com/shashank-sharma/backend/internal/util"
)

type MailAuthData struct {
	Code     string `json:"code"`
	Provider string `json:"provider"`
}

// MailAuthHandler initiates the OAuth flow for Gmail
func MailAuthHandler(ms *mail.MailService, e *core.RequestEvent) error {
	return e.JSON(http.StatusOK, map[string]interface{}{
		"url": ms.GetAuthUrl(),
	})
}

// MailAuthCallback handles the OAuth callback from Gmail
func MailAuthCallback(ms *mail.MailService, e *core.RequestEvent) error {
	pbToken := e.Request.Header.Get("Authorization")
	userId, err := util.GetUserId(pbToken)
	if err != nil {
		return e.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "Invalid authorization",
		})
	}

	mailAuthData := &MailAuthData{}
	if err := e.BindBody(mailAuthData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	googleConfig := ms.GetConfig()
	token, err := googleConfig.Exchange(context.Background(), mailAuthData.Code)
	if err != nil {
		return e.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "Invalid token exchange",
		})
	}

	logger.Debug.Println("Found token: ", token)

	client := googleConfig.Client(context.Background(), token)
	userInfo, err := oauth.FetchUserInfo(client)
	if err != nil {
		return e.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "Failed to fetch userinfo",
		})
	}

	expiry := types.DateTime{}
	expiry.Scan(token.Expiry)

	mailToken := &models.Token{
		User:         userId,
		Provider:     "gmail",
		Account:      userInfo.Email,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       expiry,
		Scope:        token.TokenType,
		IsActive:     true,
	}

	logger.Debug.Println("Set refresh token: ", token.RefreshToken)

	if err := query.UpsertRecord[*models.Token](mailToken, map[string]interface{}{
		"provider": "gmail",
		"account":  userInfo.Email,
	}); err != nil {
		return e.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "Error saving token",
		})
	}

	logger.Debug.Println("Token id: ", mailToken.Id)

	_, err = ms.InitializeLabels(mailToken.Id, userId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Error initializing labels: " + err.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "Mail authentication successful",
	})
}

// MailSyncHandler triggers a non-blocking mail sync
func MailSyncHandler(ms *mail.MailService, e *core.RequestEvent) error {
	pbToken := e.Request.Header.Get("Authorization")
	userId, err := util.GetUserId(pbToken)
	if err != nil {
		return e.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "Invalid authorization",
		})
	}

	// Find active mail sync configuration
	// TODO: Assumption that mailsync is possible by only 1 provider
	mailSync, err := query.FindByFilter[*models.MailSync](map[string]interface{}{
		"user":      userId,
		"is_active": true,
	})
	if err != nil {
		return e.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "No active mail sync configuration found",
		})
	}

	syncStatus := map[string]interface{}{
		"sync_status": "in_progress",
		"last_synced": types.NowDateTime(),
	}

	if err := query.UpdateRecord[*models.MailSync](mailSync.Id, syncStatus); err != nil {
		logger.LogError("Failed to update sync status", err)
	}

	go func() {
		logger.LogInfo("Starting async mail sync for user: " + userId)

		err := ms.SyncMessages(mailSync)

		var finalStatus string
		if err != nil {
			logger.LogError("Mail sync failed", err)
			finalStatus = "failed"
		} else {
			finalStatus = "completed"
		}

		if err := query.UpdateRecord[*models.MailSync](mailSync.Id, map[string]interface{}{
			"sync_status": finalStatus,
		}); err != nil {
			logger.LogError("Failed to update final sync status", err)
		}
	}()

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "Mail sync started in background",
		"status":  "in_progress",
	})
}

// MailSyncStatusHandler checks the status of the mail sync
func MailSyncStatusHandler(ms *mail.MailService, e *core.RequestEvent) error {
	pbToken := e.Request.Header.Get("Authorization")
	userId, err := util.GetUserId(pbToken)
	if err != nil {
		return e.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "Invalid authorization",
		})
	}

	mailSync, err := query.FindByFilter[*models.MailSync](map[string]interface{}{
		"user":      userId,
		"is_active": true,
	})
	if err != nil {
		return e.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "No active mail sync configuration found",
		})
	}

	messageCount, err := query.CountRecords[*models.MailMessage](map[string]interface{}{
		"user":      userId,
		"mail_sync": mailSync.Id,
	})
	if err != nil {
		logger.LogError("Failed to count messages", err)
		messageCount = 0
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"id":            mailSync.Id,
		"status":        mailSync.SyncStatus,
		"last_synced":   mailSync.LastSynced,
		"message_count": messageCount,
	})
}
