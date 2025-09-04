package model

type UserPermission struct {
	BaseModel
	UserID                   string `json:"user_id"` // Foreign Key to the User table
	PeoplePermission         bool   `json:"people_permission"`
	DevicePermission         bool   `json:"device_permission"`
	RulePermission           bool   `json:"rule_permission"`
	TimeAttendancePermission bool   `json:"time_attendance_permission"`
	ReportPermission         bool   `json:"report_permission"`
	NotificationPermission   bool   `json:"notification_permission"`
	SystemLogPermission      bool   `json:"system_log_permission"`
}
