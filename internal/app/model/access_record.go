package model

import "time"

type AccessRecord struct {
	PersonId              string    `json:"person_id"`
	AccessControlDeviceID string    `json:"access_control_device_id"`
	Type                  string    `json:"type"`
	Result                string    `json:"result"`
	Timestamp             string    `json:"timestamp"`
	AccessTime            time.Time `json:"access_time"`
}
