package model

type Device struct {
	BaseModel
	Name        string `json:"name"`
	HostAddress string `json:"host_address"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}
