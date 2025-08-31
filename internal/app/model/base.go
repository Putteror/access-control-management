package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel defines the common fields with a UUID ID
type BaseModel struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
