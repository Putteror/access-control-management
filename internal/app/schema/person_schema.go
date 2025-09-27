package schema

import (
	"time"
)

type PersonSearchQuery struct {
	FirstName    string `form:"firstName"`
	LastName     string `form:"lastName"`
	Company      string `form:"company"`
	Department   string `form:"department"`
	JobPosition  string `form:"jobPosition"`
	MobileNumber string `form:"mobileNumber"`
	Email        string `form:"email"`
	Page         int    `form:"page"`
	Limit        int    `form:"limit"`
	All          bool   `form:"all"`
}

var PERSON_TYPE_LIST = []string{"employee", "visitor"}

type PersonRequest struct {
	FirstName           *string  `form:"firstName" validate:"required"`
	MiddleName          *string  `form:"middleName"`
	LastName            *string  `form:"lastName" validate:"required"`
	PersonType          *string  `form:"personType" validate:"required,personType"`
	PersonID            *string  `form:"personId"`
	Gender              *string  `form:"gender"`
	DateOfBirth         *string  `form:"dateOfBirth"`
	Company             *string  `form:"company"`
	Department          *string  `form:"department"`
	JobPosition         *string  `form:"jobPosition"`
	Address             *string  `form:"address"`
	MobileNumber        *string  `form:"mobileNumber"`
	Email               *string  `form:"email"`
	IsVerified          *bool    `form:"isVerified"`
	ActiveAt            *string  `form:"activeAt"`
	ExpireAt            *string  `form:"expireAt"`
	CardIDs             []string `form:"cardIds"`
	LicensePlateTexts   []string `form:"licensePlateTexts"`
	AccessControlRuleID *string  `form:"accessControlRuleId"`
	TimeAttendanceID    *string  `form:"timeAttendanceId"`
	// Face image will receive in function
}

type PersonInfoResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PersonType string `json:"personType"`
	PersonID   string `json:"personId"`
}

type PersonResponse struct {
	ID                string                         `json:"id"`
	FirstName         string                         `json:"firstName"`
	MiddleName        *string                        `json:"middleName"`
	LastName          string                         `json:"lastName"`
	PersonType        string                         `json:"personType"`
	PersonID          *string                        `json:"personId"`
	Gender            *string                        `json:"gender"`
	DateOfBirth       *time.Time                     `json:"dateOfBirth"`
	Company           *string                        `json:"company"`
	Department        *string                        `json:"department"`
	JobPosition       *string                        `json:"jobPosition"`
	Address           *string                        `json:"address"`
	MobileNumber      *string                        `json:"mobileNumber"`
	Email             *string                        `json:"email"`
	IsVerified        bool                           `json:"isVerified"`
	CardIDs           []string                       `json:"cardIds"`
	LicensePlateTexts []string                       `json:"licensePlateTexts"`
	FaceImagePath     *string                        `json:"faceImagePath"`
	ActiveAt          *time.Time                     `json:"activeAt"`
	ExpireAt          *time.Time                     `json:"expireAt"`
	AccessControlRule *AccessControlRuleInfoResponse `json:"accessControlRule"`
	TimeAttendance    *AttendanceInfoResponse        `json:"timeAttendance"`
}
