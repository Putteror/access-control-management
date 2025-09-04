package model

type RegisterFormFieldAnswer struct {
	BaseModel
	Name                string `json:"name"`
	RegisterFormID      string `json:"register_form_id"`
	RegisterFormFieldID string `json:"register_form_field_id"`
	PersonID            string `json:"person_id"`
	Answer              string `json:"answer"`
}
