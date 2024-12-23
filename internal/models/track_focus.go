package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ models.Model = (*TrackFocus)(nil)

type TrackFocus struct {
	models.BaseModel

	User      string                  `db:"user" json:"user"`
	Device    string                  `db:"device" json:"device"`
	Tags      types.JsonArray[string] `db:"tags" json:"tags"`
	Metadata  string                  `db:"metadata" json:"metadata"`
	BeginDate types.DateTime          `db:"begin_date" json:"begin_date"`
	EndDate   types.DateTime          `db:"end_date" json:"end_date"`
}

func (m *TrackFocus) TableName() string {
	return "track_focus"
}
