package models

import (
	pbModels "github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
)

type Model interface {
	TableName() string
	PK() any
	LastSavedPK() any
	IsNew() bool
	SetId(id string)
	GetId() string
	MarkAsNew()
	MarkAsNotNew()
}

type BaseModel struct {
	Model
	lastSavedPK string

	Id string `db:"id" json:"id" form:"id" xml:"id"`
}

// HasId returns whether the model has a nonzero id.
func (m *BaseModel) HasId() bool {
	return m.GetId() != ""
}

// GetId returns the model id.
func (m *BaseModel) GetId() string {
	return m.Id
}

// SetId sets the model id to the provided string value.
func (m *BaseModel) SetId(id string) {
	m.Id = id
}

// RefreshId generates and sets a new model id.
//
// The generated id is a cryptographically random 15 characters length string.
func (m *BaseModel) RefreshId() {
	m.Id = security.RandomStringWithAlphabet(pbModels.DefaultIdLength, pbModels.DefaultIdAlphabet)
}
