package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

// ensures that the Article struct satisfy the models.Model interface
var _ models.Model = (*TrackItems)(nil)

type TrackItems struct {
	models.BaseModel

	User      string         `db:"user" json:"user"`
	TrackId   int64          `db:"track_id" json:"track_id"`
	Source    string         `db:"source" json:"source"`
	App       string         `db:"app" json:"app"`
	TaskName  string         `db:"task_name" json:"task_name"`
	Title     string         `db:"title" json:"title"`
	BeginDate types.DateTime `db:"begin_date" json:"begin_date"`
	EndDate   types.DateTime `db:"end_date" json:"end_date"`
}

func (m *TrackItems) TableName() string {
	return "track_items" // the name of your collection
}
