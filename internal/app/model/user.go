package model

type User struct {
	BaseModel
	Username     string `json:"username" gorm:"unique"`
	PasswordHash string `json:"-"`
	PermissionID string `json:"permission_id"` // Foreign Key to UserPermission table
}
