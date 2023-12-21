package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

type TrackUploadAPI struct {
	Source     string           `json:"source" form:"source"`
	ForceCheck bool             `json:"force_check" form:"force_check"`
	File       *filesystem.File `json:"file" form:"file"`
}

// ensures that the Article struct satisfy the models.Model interface
var _ models.Model = (*TrackUpload)(nil)

type TrackUpload struct {
	models.BaseModel

	User            string `db:"user" json:"user"`
	Source          string `db:"source" json:"source"`
	File            string `db:"file" json:"file"`
	Synced          bool   `db:"synced" json:"synced"`
	Status          string `db:"status" json:"status"`
	TotalRecord     int64  `db:"total_record" json:"total_record"`
	DuplicateRecord int64  `db:"duplicate_record" json:"duplicate_record"`
}

func (m *TrackUpload) TableName() string {
	return "track_upload" // the name of your collection
}
