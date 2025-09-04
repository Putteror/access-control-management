package model

type RegisterFormField struct {
	BaseModel
	Name           string `json:"name"`
	RegisterFormID string `json:"register_form_id"`
	FieldType      string `json:"field_type"`
	InputType      string `json:"input_type"`
	Placeholder    string `json:"placeholder"`
	Label          string `json:"label"`
	HelpText       string `json:"help_text"`
	IsRequired     bool   `json:"is_required"`
	DefaultValue   string `json:"default_value"`
}
