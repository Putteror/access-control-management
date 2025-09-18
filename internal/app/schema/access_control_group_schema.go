package schema

// AccessControlDeviceSearchQuery defines the search parameters for devices.
type AccessControlGroupSearchQuery struct {
	Name  string `json:"name"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

// Request

type AccessControlGroupScheduleRequest struct {
	DayOfWeek int     `json:"dayOfWeek"`
	Date      *string `json:"date"`
	StartTime string  `json:"startTime"`
	EndTime   string  `json:"endTime"`
}

type AccessControlGroupRequest struct {
	Name                        *string                             `json:"name" validate:"required"`
	AccessControlDeviceIDs      []string                            `json:"accessControlDeviceIds"`
	AccessControlGroupSchedules []AccessControlGroupScheduleRequest `json:"accessControlSchedules"`
}

// Response

type AccessControlGroupScheduleResponse struct {
	DayOfWeek int     `json:"dayOfWeek"`
	Date      *string `json:"date"`
	StartTime string  `json:"startTime"`
	EndTime   string  `json:"endTime"`
}

type AccessControlGroupInfoResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AccessControlGroupResponse struct {
	ID                          string                               `json:"id"`
	Name                        string                               `json:"name"`
	AccessControlDevices        []AccessControlDeviceInfoResponse    `json:"accessControlDevices"`
	AccessControlGroupSchedules []AccessControlGroupScheduleResponse `json:"accessControlSchedules"`
}
