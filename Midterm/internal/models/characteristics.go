package models

type Characteristics struct {
	ID       int       `json:"id"`
	Property *Property `json:"property"`
	Value    string    `json:"value"`
}
