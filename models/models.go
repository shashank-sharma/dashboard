package models

import (
	pbModels "github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Model interface {
	TableName() string
	IsNew() bool
	MarkAsNew()
	MarkAsNotNew()
	HasId() bool
	GetId() string
	SetId(id string)
	GetCreated() types.DateTime
	GetUpdated() types.DateTime
	RefreshId()
	RefreshCreated()
	RefreshUpdated()
}

type BaseModel struct {
	Model
	isNewFlag bool

	Id      string         `db:"id" json:"id"`
	Created types.DateTime `db:"created" json:"created"`
	Updated types.DateTime `db:"updated" json:"updated"`
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

// MarkAsNew sets the model isNewFlag enforcing [m.IsNew()] to be true.
func (m *BaseModel) MarkAsNew() {
	m.isNewFlag = true
}

// UnmarkAsNew resets the model isNewFlag.
func (m *BaseModel) MarkAsNotNew() {
	m.isNewFlag = false
}

// IsNew indicates what type of db query (insert or update)
// should be used with the model instance.
func (m *BaseModel) IsNew() bool {
	return m.isNewFlag || !m.HasId()
}

// GetCreated returns the model Created datetime.
func (m *BaseModel) GetCreated() types.DateTime {
	return m.Created
}

// GetUpdated returns the model Updated datetime.
func (m *BaseModel) GetUpdated() types.DateTime {
	return m.Updated
}

// RefreshId generates and sets a new model id.
//
// The generated id is a cryptographically random 15 characters length string.
func (m *BaseModel) RefreshId() {
	m.Id = security.RandomStringWithAlphabet(pbModels.DefaultIdLength, pbModels.DefaultIdAlphabet)
}

// RefreshCreated updates the model Created field with the current datetime.
func (m *BaseModel) RefreshCreated() {
	m.Created = types.NowDateTime()
}

// RefreshUpdated updates the model Updated field with the current datetime.
func (m *BaseModel) RefreshUpdated() {
	m.Updated = types.NowDateTime()
}
