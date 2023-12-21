package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

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

// ensures that the Article struct satisfy the models.Model interface
var _ models.Model = (*TrackDevice)(nil)

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

func (m *TrackDevice) TableName() string {
	return "devices" // the name of your collection
}
