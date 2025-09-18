package model

import "time"

type AccessRecord struct {
	PersonID              string    `json:"person_id"`
	AccessControlDeviceID string    `json:"access_control_device_id"`
	Type                  string    `json:"type"`
	Result                string    `json:"result"`
	AccessTime            time.Time `json:"access_time"`
}
