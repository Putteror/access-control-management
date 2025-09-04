package model

type AttendanceRecord struct {
	BaseModel
	PersonId             string `json:"person_id"`
	AttendanceScheduleID string `json:"attendance_schedule_id"`
	AccessRecordId       string `json:"access_record_id"`
	Date                 string `json:"date"`
	Status               string `json:"status"`
}
