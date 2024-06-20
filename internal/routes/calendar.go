package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
	"github.com/shashank-sharma/backend/internal/services/calendar"
	"github.com/shashank-sharma/backend/internal/util"
)

func CalendarSyncHandler(cs *calendar.CalendarService, c echo.Context) error {
	pbToken := c.Request().Header.Get(echo.HeaderAuthorization)
	userId, _ := util.GetUserId(pbToken)
	calendarSync, err := query.FindByFilter[*models.CalendarSync](map[string]interface{}{
		"user":      userId,
		"is_active": true,
	})

	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Calendar sync not found"})
	}

	if err := cs.SyncEvents(calendarSync); err != nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Failed to sync"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Sync completed"})
}
