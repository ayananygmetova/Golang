package models

type Product struct {
	ID           int     `json:"id" db:"id"`
	Name         string  `json:"name" db:"name"`
	Manufacturer string  `json:"manufacturer" db:"manufacturer"`
	Description  string  `json:"description" db:"description"`
	Price        float32 `json:"price" db:"price"`
	Brand        string  `json:"brand" db:"brand"`
	CategoryId   int     `json:"category_id"  db:"category_id"`
}
