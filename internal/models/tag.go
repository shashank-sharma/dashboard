package models

import (
	"github.com/pocketbase/pocketbase/core"
)

var _ core.Model = (*Tag)(nil)

type Tag struct {
	BaseModel

	User        string `db:"user" json:"user"`
	Name        string `db:"name" json:"name"`
	Color       string `db:"color" json:"color"`
	Description string `db:"description" json:"description"`
	IsAICreated bool   `db:"is_ai_created" json:"is_ai_created"`
}

func (m *Tag) TableName() string {
	return "tags"
} 