package models

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ core.Model = (*Token)(nil)

type Token struct {
	BaseModel

	User         string         `db:"user" json:"user"`
	Provider     string         `db:"provider" json:"provider"`
	Account      string         `db:"account" json:"account"`
	AccessToken  string         `db:"access_token" json:"access_token"`
	TokenType    string         `db:"token_type" json:"token_type"`
	RefreshToken string         `db:"refresh_token" json:"refresh_token"`
	Expiry       types.DateTime `db:"expiry" json:"expiry"`
	Scope        string         `db:"scope" json:"scope"`
	IsActive     bool           `db:"is_active" json:"is_active"`
	LastUsed     types.DateTime `db:"last_used" json:"last_used"`
}

func (m *Token) TableName() string {
	return "tokens"
}
