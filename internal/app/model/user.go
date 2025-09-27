package model

type User struct {
	BaseModel
	Username     string         `json:"username" gorm:"unique"`
	PasswordHash string         `json:"-"`
	PermissionID string         // Foreign Key to UserPermission table
	Permission   UserPermission `json:"permission" gorm:"foreignKey:PermissionID"`
	Status       string         `json:"status"`
}
