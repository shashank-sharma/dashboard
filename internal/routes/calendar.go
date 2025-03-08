package routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
	"github.com/shashank-sharma/backend/internal/services/calendar"
	"github.com/shashank-sharma/backend/internal/services/oauth"
	"github.com/shashank-sharma/backend/internal/util"
)

type CalendarTokenAPI struct {
	Code     string `json:"code"`
	Provider string `json:"provider"`
}

func RegisterCalendarRoutes(apiRouter *router.RouterGroup[*core.RequestEvent], path string, calendarService *calendar.CalendarService) {
	calendarRouter := apiRouter.Group(path)
	calendarRouter.GET("/auth/redirect", func(e *core.RequestEvent) error {
		return CalendarAuthHandler(calendarService, e)
	})
	calendarRouter.POST("/auth/callback", func(e *core.RequestEvent) error {
		return CalendarAuthCallback(calendarService, e)
	})
	calendarRouter.POST("/sync", func(e *core.RequestEvent) error {
		return CalendarSyncHandler(calendarService, e)
	})
}

func CalendarAuthHandler(cs *calendar.CalendarService, e *core.RequestEvent) error {
	return e.JSON(http.StatusOK, map[string]interface{}{"url": cs.GetAuthUrl()})
}

func CalendarAuthCallback(cs *calendar.CalendarService, e *core.RequestEvent) error {
	pbToken := e.Request.Header.Get("Authorization")
	userId, err := util.GetUserId(pbToken)

	googleConfig := cs.GetConfig()

	calTokenData := &CalendarTokenAPI{}
	if err := e.BindBody(calTokenData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	token, err := googleConfig.Exchange(context.Background(), calTokenData.Code)
	if err != nil {
		return e.JSON(http.StatusForbidden, map[string]interface{}{"message": "Invalid token exchange"})
	}

	client := googleConfig.Client(context.Background(), token)
	userInfo, err := oauth.FetchUserInfo(client)
	if err != nil {
		return e.JSON(http.StatusForbidden, map[string]interface{}{"message": "Failed to fetch userinfo"})
	}

	expiry := types.DateTime{}
	expiry.Scan(token.Expiry)

	calToken := &models.CalendarToken{
		User:         userId,
		Account:      userInfo.Email,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       expiry,
	}

	if err := query.UpsertRecord[*models.CalendarToken](calToken, map[string]interface{}{"account": userInfo.Email}); err != nil {
		return e.JSON(http.StatusForbidden, map[string]interface{}{"message": "Error saving record"})
	}
	return e.JSON(http.StatusOK, map[string]interface{}{"message": "Authenticated successfully"})
}


func CalendarSyncHandler(cs *calendar.CalendarService, e *core.RequestEvent) error {
	pbToken := e.Request.Header.Get("Authorization")
	userId, _ := util.GetUserId(pbToken)
	calendarSync, err := query.FindByFilter[*models.CalendarSync](map[string]interface{}{
		"user":      userId,
		"is_active": true,
	})

	if err != nil {
		return e.JSON(http.StatusForbidden, map[string]interface{}{"message": "Calendar sync not found"})
	}

	if err := cs.SyncEvents(calendarSync); err != nil {
		return e.JSON(http.StatusForbidden, map[string]interface{}{"message": "Failed to sync"})
	}
	return e.JSON(http.StatusOK, map[string]interface{}{"message": "Sync completed"})
}
