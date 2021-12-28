package models

type Characteristics struct {
	ID         int    `json:"id" db:"id"`
	PropertyId string `json:"property_id" db:"property_id"`
	Value      string `json:"value" db:"value"`
}
