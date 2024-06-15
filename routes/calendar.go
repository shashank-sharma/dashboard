package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/shashank-sharma/backend/models"
	"github.com/shashank-sharma/backend/query"
	"github.com/shashank-sharma/backend/services/calendar"
	"github.com/shashank-sharma/backend/util"
)

func CalendarSyncHandler(c echo.Context) error {
	pbToken := c.Request().Header.Get(echo.HeaderAuthorization)
	userId, _ := util.GetUserId(pbToken)
	calendarSync, err := query.FindByFilter[*models.CalendarSync](map[string]interface{}{
		"user":      userId,
		"is_active": true,
	})

	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Calendar sync not found"})
	}

	if err := calendar.SyncEvents(calendarSync); err != nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Failed to sync"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Sync completed"})
}
