package model

type PersonLicensePlate struct {
	BaseModel
	LicensePlateText string `json:"license_plate_text" gorm:"unique"`
	PersonID         string `json:"person_id"` // FK to people table
}
