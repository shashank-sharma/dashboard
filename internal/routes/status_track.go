package routes

import (
	"net/http"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
)

func TrackDeviceStatus(e *core.RequestEvent) error {
	data := &models.TrackDeviceUpdateAPI{}

	if err := e.BindBody(data); err != nil || data.UserId == "" || data.Token == "" || data.ProductId == "" {
		// TODO: Simply say unauthorized
		return apis.NewBadRequestError("Failed to read request data", err)
	}

	record, err := query.FindByFilter[*models.DevToken](map[string]interface{}{
		"user":  data.UserId,
		"token": data.Token,
	})

	if err != nil || record == nil {
		return e.JSON(http.StatusForbidden, map[string]interface{}{"message": "Not authorized"})
	}

	deviceRecord, err := query.FindByFilter[*models.TrackDevice](map[string]interface{}{
		"user": data.UserId,
		"id":   data.ProductId,
	})

	if err != nil || deviceRecord == nil {
		return e.JSON(http.StatusForbidden, map[string]interface{}{"message": "Authorized, but product not found"})
	}

	err = query.UpdateRecord[*models.TrackDevice](deviceRecord.Id, map[string]interface{}{
		"is_active": true,
		"is_online": true,
	})

	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]interface{}{"message": err})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{"message": "ok"})
}
