package models

var _ Model = (*Users)(nil)

type Users struct {
	BaseModel

	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Name     string `db:"name" json:"name"`
}

func (m *Users) TableName() string {
	return "users"
}
