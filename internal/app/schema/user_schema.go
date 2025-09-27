package schema

type UserSearchQuery struct {
	Username string `form:"username"`
	Status   string `form:"status"`
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
}

type UserPermissionRequest struct {
	PeoplePermission         *bool `json:"peoplePermission"`
	DevicePermission         *bool `json:"devicePermission"`
	RulePermission           *bool `json:"rulePermission"`
	TimeAttendancePermission *bool `json:"timeAttendancePermission"`
	ReportPermission         *bool `json:"reportPermission"`
	NotificationPermission   *bool `json:"notificationPermission"`
	SystemLogPermission      *bool `json:"systemLogPermission"`
}

type UserRequest struct {
	Username   *string                `json:"username"`
	Password   *string                `json:"password"`
	Status     *string                `json:"status"`
	Permission *UserPermissionRequest `json:"permission" `
}

type UserPermissionResponse struct {
	ID                       string `json:"id"`
	PeoplePermission         bool   `json:"peoplePermission"`
	DevicePermission         bool   `json:"devicePermission"`
	RulePermission           bool   `json:"rulePermission"`
	TimeAttendancePermission bool   `json:"timeAttendancePermission"`
	ReportPermission         bool   `json:"reportPermission"`
	NotificationPermission   bool   `json:"notificationPermission"`
	SystemLogPermission      bool   `json:"systemLogPermission"`
}

type UserResponse struct {
	ID         string                 `json:"id"`
	Username   string                 `json:"username"`
	Status     string                 `json:"status"`
	Permission UserPermissionResponse `json:"permission"`
}
