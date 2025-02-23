package models

import (
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ Model = (*Notification)(nil)
var _ Model = (*NotificationSettings)(nil)

type Notification struct {
	BaseModel

	User     string         `db:"user" json:"user"`
	Type     string         `db:"type" json:"type"`
	Title    string         `db:"title" json:"title"`
	Content  string         `db:"content" json:"content"`
	Priority string         `db:"priority" json:"priority"`
	Status   string         `db:"status" json:"status"`
	Metadata string         `db:"metadata" json:"metadata"`
	ReadAt   types.DateTime `db:"read_at" json:"read_at"`
}

type NotificationSettings struct {
	BaseModel

	User            string                  `db:"user" json:"user"`
	EmailEnabled    bool                    `db:"email_enabled" json:"email_enabled"`
	TypesEnabled    types.JSONArray[string] `db:"types_enabled" json:"types_enabled"`
	QuietHoursStart string                  `db:"quiet_hours_start" json:"quiet_hours_start"`
	QuietHoursEnd   string                  `db:"quiet_hours_end" json:"quiet_hours_end"`
}

func (m *Notification) TableName() string {
	return "notifications"
}

func (m *NotificationSettings) TableName() string {
	return "notification_settings"
}
