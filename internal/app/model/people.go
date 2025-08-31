package model

type People struct {
	BaseModel
	Name     string `json:"name"`
	Position string `json:"position"`
	Admin    bool   `json:"admin"`
}
