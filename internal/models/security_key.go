package models


type SecurityKey struct {
	BaseModel
	User        string         `json:"user"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	PrivateKey  string         `json:"private_key"`
	PublicKey   string         `json:"public_key"`
	IsActive    bool           `json:"is_active"`
}

func (s *SecurityKey) TableName() string {
	return "security_keys"
} 