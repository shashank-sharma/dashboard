package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ models.Model = (*FoldToken)(nil)

type FoldToken struct {
	models.BaseModel

	User           string         `db:"user" json:"user"`
	Phone          string         `db:"phone" json:"phone"`
	Uuid           string         `db:"uuid" json:"uuid"`
	FirstName      string         `db:"first_name" json:"first_name"`
	LastName       string         `db:"last_name" json:"last_name"`
	Email          string         `db:"email" json:"email"`
	AccessToken    string         `db:"access_token" json:"access_token"`
	RefreshToken   string         `db:"refresh_token" json:"refresh_token"`
	UserAgent      string         `db:"user_agent" json:"user_agent"`
	DeviceHash     string         `db:"device_hash" json:"device_hash"`
	DeviceLocation string         `db:"device_location" json:"device_location"`
	DeviceName     string         `db:"device_name" json:"device_name"`
	DeviceType     string         `db:"device_type" json:"device_type"`
	ExpiresAt      types.DateTime `db:"expires_at" json:"expires_at"`
}

func (m *FoldToken) TableName() string {
	return "fold_tokens"
}
