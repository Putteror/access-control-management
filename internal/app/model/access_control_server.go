package model

import (
	"time"
)

type AccessControlServer struct {
	BaseModel
	Name        string     `json:"name"`
	HostAddress string     `json:"host_address"`
	Username    string     `json:"username"`
	Password    string     `json:"password"`
	AccessToken string     `json:"access_token"`
	ApiToken    string     `json:"api_token"`
	Status      string     `json:"status"`
	LastSyncAt  *time.Time `json:"last_sync_at"`
}
