package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
)

func TrackDeviceStatus(c echo.Context) error {
	logger.Debug.Println("Started track")
	data := &models.TrackDeviceUpdateAPI{}

	if err := c.Bind(data); err != nil || data.UserId == "" || data.Token == "" || data.ProductId == "" {
		logger.Debug.Println("Error in parsing =", err)
		// TODO: Simply say unauthorized
		return apis.NewBadRequestError("Failed to read request data", err)
	}

	record, err := query.FindByFilter[*models.DevToken](map[string]interface{}{
		"user":  data.UserId,
		"token": data.Token,
	})

	logger.Debug.Println("Found token record: ", record)

	if err != nil || record == nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Not authorized"})
	}

	deviceRecord, err := query.FindByFilter[*models.TrackDevice](map[string]interface{}{
		"user": data.UserId,
		"id":   data.ProductId,
	})

	if err != nil || deviceRecord == nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Authorized, but product not found"})
	}

	err = query.UpdateRecord[*models.TrackDevice](deviceRecord.Id, map[string]interface{}{
		"is_active": true,
		"is_online": true,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": err})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "ok"})
}
