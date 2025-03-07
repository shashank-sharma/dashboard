package models

type Server struct {
	BaseModel
	User        string         `json:"user"`
	Name        string         `json:"name"`
	Provider    string         `json:"provider"`
	IP          string         `json:"ip"`
	Port        int            `json:"port"`
	Username    string         `json:"username"`
	SecurityKey string         `json:"security_key"`
	SSHEnabled  bool           `json:"ssh_enabled"`
	IsActive    bool           `json:"is_active"`
	IsReachable bool           `json:"is_reachable"`
}

func (s *Server) TableName() string {
	return "servers"
} 