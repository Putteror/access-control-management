package model

import "time"

type PersonCard struct {
	BaseModel
	CardNumber string    `json:"card_number" gorm:"unique"`
	PersonID   string    `json:"person_id"` // FK to people table
	ActiveAt   time.Time `json:"active_at"`
	ExpireAt   time.Time `json:"expire_at"`
}
