package model

import (
	"time"
)

type TimeRecord struct {
	BaseModel
	PeopleID     uint      `json:"people_id"`
	ClockInTime  time.Time `json:"clock_in_time"`
	ClockOutTime time.Time `json:"clock_out_time"`
}
