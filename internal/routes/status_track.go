package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/shashank-sharma/backend/internal/config"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
)

func TrackDeviceStatus(c echo.Context) error {
	app := config.GetApp()
	logger.Debug.Println("Started track")
	data := &models.TrackDeviceUpdateAPI{}

	if err := c.Bind(data); err != nil || data.UserId == "" || data.Token == "" || data.ProductId == "" {
		logger.Debug.Println("Error in parsing =", err)
		// TODO: Simply say unauthorized
		return apis.NewBadRequestError("Failed to read request data", err)
	}
	record, err := app.Dao().FindFirstRecordByFilter(
		"dev_tokens", "user = {:user} && token = {:token}",
		dbx.Params{"user": data.UserId},
		dbx.Params{"token": data.Token},
	)

	if err != nil || record == nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Not authorized"})
	}

	deviceRecord, err := app.Dao().FindFirstRecordByFilter(
		"devices", "user = {:user} && id = {:id}",
		dbx.Params{"user": data.UserId},
		dbx.Params{"id": data.ProductId},
	)

	if err != nil || deviceRecord == nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Authorized, but product not found"})
	}

	deviceRecord.Set("is_active", true)
	deviceRecord.Set("is_online", true)

	if err := app.Dao().SaveRecord(deviceRecord); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": err})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "ok"})
}
