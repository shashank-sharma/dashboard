package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ models.Model = (*TrackDevice)(nil)
var _ models.Model = (*TrackItems)(nil)
var _ models.Model = (*TrackUpload)(nil)

// Use structure embedding
type TrackDeviceAPI struct {
	Name     string `json:"name" form:"name"`
	HostName string `json:"hostname" form:"hostname"`
	Os       string `db:"os" json:"os"`
	Arch     string `db:"arch" json:"arch"`
}

type TrackDeviceUpdateAPI struct {
	UserId    string `json:"userid" json:"userid"`
	Token     string `json:"token" json:"token"`
	ProductId string `json:"productid" json:"productid"`
}

type TrackDevice struct {
	models.BaseModel

	User      string         `db:"user" json:"user"`
	Name      string         `db:"name" json:"name"`
	HostName  string         `db:"hostname" json:"hostname"`
	Os        string         `db:"os" json:"os"`
	Arch      string         `db:"arch" json:"arch"`
	IsOnline  bool           `db:"is_online" json:"is_online"`
	IsActive  bool           `db:"is_active" json:"is_active"`
	BeginDate types.DateTime `db:"begin_date" json:"begin_date"`
	EndDate   types.DateTime `db:"end_date" json:"end_date"`
}

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

type TrackUploadAPI struct {
	Source     string           `json:"source" form:"source"`
	ForceCheck bool             `json:"force_check" form:"force_check"`
	File       *filesystem.File `json:"file" form:"file"`
}

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

func (m *TrackDevice) TableName() string {
	return "devices"
}

func (m *TrackItems) TableName() string {
	return "track_items"
}

func (m *TrackUpload) TableName() string {
	return "track_upload"
}
