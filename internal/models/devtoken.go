package models

import (
	"github.com/pocketbase/pocketbase/models"
)

var _ models.Model = (*DevToken)(nil)

type DevToken struct {
	models.BaseModel

	User     string `db:"user" json:"user"`
	Token    string `db:"token" json:"token"`
	IsActive bool   `db:"is_active" json:"is_active"`
}

func (m *DevToken) TableName() string {
	return "dev_tokens"
}
