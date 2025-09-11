package dto

import "time"

// PersonDTO represents the data to be returned for a person.
// It includes only the necessary fields and uses JSON tags to format the output.
type PersonDTO struct {
	ID           string     `json:"id"`
	FirstName    string     `json:"firstName"`
	LastName     string     `json:"lastName"`
	Email        string     `json:"email"`
	DateOfBirth  *time.Time `json:"dateOfBirth,omitempty"` // Example of omitempty
	MobileNumber string     `json:"mobileNumber,omitempty"`
}
