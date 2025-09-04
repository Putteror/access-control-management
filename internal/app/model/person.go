package model

import "time"

type Person struct {
	BaseModel
	FirstName        string    `json:"first_name"`
	MiddleName       *string   `json:"middle_name"`
	LastName         string    `json:"last_name"`
	PersonType       string    `json:"person_type"`
	Company          string    `json:"company"`
	Department       string    `json:"department"`
	JobPosition      string    `json:"job_position"`
	Address          string    `json:"address"`
	MobileNumber     string    `json:"mobile_number"`
	Email            string    `json:"email"`
	FaceImagePath    string    `json:"face_image_path"`
	IsVerified       bool      `json:"is_verified" gorm:"default:false"`
	ActiveAt         time.Time `json:"active_at"`
	ExpireAt         time.Time `json:"expire_at"`
	RuleID           string    `json:"rule_id"`
	TimeAttendanceID string    `json:"time_attendance_id"`
}
