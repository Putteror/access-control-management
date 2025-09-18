package schema

// AccessControlDeviceSearchQuery defines the search parameters for devices.
type AccessControlDeviceSearchQuery struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	HostAddress string `json:"hostAddress"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
}

type AccessControlDeviceRequest struct {
	Name                  *string `json:"name" validate:"required"`
	Type                  *string `json:"type" validate:"required"`
	HostAddress           *string `json:"hostAddress" validate:"required"`
	Username              *string `json:"username"`
	Password              *string `json:"password"`
	AccessToken           *string `json:"accessToken"`
	ApiToken              *string `json:"apiToken"`
	RecordScan            *bool   `json:"recordScan"`
	RecordAttendance      *bool   `json:"recordAttendance"`
	AllowClockIn          *bool   `json:"allowClockIn"`
	AllowClockOut         *bool   `json:"allowClockOut"`
	Status                *string `json:"status"`
	AccessControlServerID *string `json:"accessControlServerId"`
}

type AccessControlDeviceInfoResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	HostAddress string `json:"hostAddress"`
}

type AccessControlDeviceResponse struct {
	ID                  string                           `json:"id"`
	Name                string                           `json:"name"`
	HostAddress         string                           `json:"hostAddress"`
	Type                string                           `json:"type"`
	Status              string                           `json:"status"`
	Username            *string                          `json:"username"`
	Password            *string                          `json:"password"`
	AccessToken         *string                          `json:"accessToken"`
	ApiToken            *string                          `json:"apiToken"`
	RecordScan          bool                             `json:"recordScan"`
	RecordAttendance    bool                             `json:"recordAttendance"`
	AllowClockIn        bool                             `json:"allowClockIn"`
	AllowClockOut       bool                             `json:"allowClockOut"`
	AccessControlServer *AccessControlServerInfoResponse `json:"accessControlServer"`
}
