package routes

import (
	"net/http"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
)

type TrackFocusAPI struct {
	User      string                  `db:"user" json:"user"`
	Device    string                  `db:"device" json:"device"`
	Tags      types.JSONArray[string] `db:"tags" json:"tags"`
	Metadata  string                  `db:"metadata" json:"metadata"`
	BeginDate types.DateTime          `db:"begin_date" json:"begin_date"`
	EndDate   types.DateTime          `db:"end_date" json:"end_date"`
}

func TrackFocus(e *core.RequestEvent) error {
	token := e.Request.Header.Get("AuthSyncToken")
	if token == "" {
		return e.JSON(http.StatusForbidden, map[string]interface{}{"message": "Dev Token missing"})
	}
	userId, err := query.ValidateDevToken(token)
	if err != nil {
		return e.JSON(http.StatusForbidden, map[string]interface{}{"message": "Failed to fetch id, token misconfigured"})
	}
	data := TrackFocusAPI{}
	if err := e.BindBody(&data); err != nil {
		logger.Error.Println("Error in parsing =", err)
		return e.JSON(http.StatusForbidden, map[string]interface{}{"message": "Failed binding data"})
	}

	_, err = query.FindByFilter[*models.TrackFocus](map[string]interface{}{
		"user":       userId,
		"device":     data.Device,
		"tags":       data.Tags,
		"metadata":   data.Metadata,
		"begin_date": data.BeginDate,
		"end_date":   data.EndDate,
	})

	if err == nil {
		logger.Error.Println("Found record, need to skip")
	}

	trackFocus := &models.TrackFocus{
		User:      userId,
		Device:    data.Device,
		Tags:      data.Tags,
		Metadata:  data.Metadata,
		BeginDate: data.BeginDate,
		EndDate:   data.EndDate,
	}

	err = query.UpsertRecord[*models.TrackFocus](trackFocus, map[string]interface{}{
		"user":       userId,
		"device":     data.Device,
		"tags":       data.Tags,
		"metadata":   data.Metadata,
		"begin_date": data.BeginDate,
	})
	if err != nil {
		logger.Error.Println("Failed updating record", err)
	}
	return e.JSON(http.StatusOK, map[string]interface{}{"message": "Created successfully"})
}
