package models

type Category struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type CategoriesFilter struct {
	Query *string `json:"query"`
}
