package model

type AttendanceSchedule struct {
	BaseModel
	AttendanceID    string  `json:"attendance_id"`
	DayOfWeek       int     `json:"day_of_week"`
	Date            *string `json:"date"`
	StartTime       string  `json:"start_time"`
	EndTime         string  `json:"end_time"`
	EarlyInMinutes  int     `json:"early_in_minutes"`
	LateInMinutes   int     `json:"late_in_minutes"`
	EarlyOutMinutes int     `json:"early_out_minutes"`
	LateOutMinutes  int     `json:"late_out_minutes"`
}
