package model

type AccessControlDevice struct {
	BaseModel
	Name                  string  `json:"name"`
	Type                  string  `json:"type"`
	HostAddress           string  `json:"host_address"`
	Username              *string `json:"username"`
	Password              *string `json:"password"`
	AccessToken           *string `json:"access_token"`
	ApiToken              *string `json:"api_token"`
	AccessControlServerID *string `json:"access_control_server_id"`
	RecordScan            bool    `json:"record_scan"`
	RecordAttendance      bool    `json:"record_attendance"`
	AllowClockIn          bool    `json:"allow_clock_in"`
	AllowClockOut         bool    `json:"allow_clock_out"`
	Status                string  `json:"status"`
}
