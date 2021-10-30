package models

type Product struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Manufacturer string  `json:"manufacturer"`
	Description  string  `json:"description"`
	Price        float32 `json:"price"`
	Brand        string  `json:"brand"`
	CategoryId   int     `json:"category_id"`
}
