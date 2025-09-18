package schema

type AttendanceSearchQuery struct {
	Name  string `form:"name"`
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
}

type AttendanceRequest struct {
	Name               *string                     `form:"name" validate:"required"`
	AttendanceSchedule []AttendanceScheduleRequest `json:"attendanceSchedules"`
}

type AttendanceScheduleRequest struct {
	DayOfWeek       *int    `json:"dayOfWeek"`
	Date            *string `json:"date"`
	StartTime       *string `json:"startTime"`
	EndTime         *string `json:"endTime"`
	EarlyInMinutes  *int    `json:"early_in_minutes"`
	LateInMinutes   *int    `json:"late_in_minutes"`
	EarlyOutMinutes *int    `json:"early_out_minutes"`
	LateOutMinutes  *int    `json:"late_out_minutes"`
}

type AttendanceScheduleResponse struct {
	ID              string  `json:"id"`
	DayOfWeek       int     `json:"dayOfWeek"`
	Date            *string `json:"date"`
	StartTime       string  `json:"startTime"`
	EndTime         string  `json:"endTime"`
	EarlyInMinutes  int     `json:"early_in_minutes"`
	LateInMinutes   int     `json:"late_in_minutes"`
	EarlyOutMinutes int     `json:"early_out_minutes"`
	LateOutMinutes  int     `json:"late_out_minutes"`
}

type AttendanceInfoResponse struct {
	ID                  string                       `json:"id"`
	Name                string                       `json:"name"`
	AttendanceSchedules []AttendanceScheduleResponse `json:"attendanceSchedules"`
}
