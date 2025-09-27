package schema

type AccessRecordSearchQuery struct {
	PersonID              string `form:"personID"`
	AccessControlDeviceID string `form:"accessControlDeviceID"`
	Type                  string `form:"type"`
	Result                string `form:"result"`
	AccessTime            string `form:"accessTime"`
	Page                  int    `form:"page"`
	Limit                 int    `form:"limit"`
}

type AccessRecordRequest struct {
	PersonID              *string `json:"personID" `
	AccessControlDeviceID *string `json:"accessControlDeviceID"`
	Type                  *string `json:"type" validate:"required"`
	Result                *string `json:"result" validate:"required"`
	AccessTime            *string `json:"accessTime" validate:"required"`
}

type AccessRecordPersonResponse struct {
	ID          string  `json:"id"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Company     *string `json:"company"`
	Department  *string `json:"department"`
	JobPosition *string `json:"jobPosition"`
}

type AccessRecordDeviceResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	HostAddress string `json:"hostAddress"`
	Type        string `json:"type"`
}

type AccessRecordResponse struct {
	ID                  string                      `json:"id"`
	Person              *AccessRecordPersonResponse `json:"person"`
	AccessControlDevice *AccessRecordDeviceResponse `json:"accessControlDevice"`
	Type                string                      `json:"type"`
	Result              string                      `json:"result"`
	AccessTime          string                      `json:"accessTime"`
}
