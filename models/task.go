package models

import "github.com/pocketbase/pocketbase/tools/types"

var _ Model = (*Task)(nil)

type Task struct {
	BaseModel

	User        string         `db:"user" json:"user"`
	Title       string         `db:"title" json:"title"`
	Description string         `db:"description" json:"description"`
	Due         types.DateTime `db:"due" json:"due"`
	Image       string         `db:"image" json:"image"`
	Category    string         `db:"category" json:"category"`
}

func (m *Task) TableName() string {
	return "tasks"
}
