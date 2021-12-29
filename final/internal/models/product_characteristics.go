package models

type ProductCharacteristics struct {
	ProductId         int `json:"product_id" db:"product_id"`
	CharacteristicsId int `json:"characteristics_id" db:"characteristics_id"`
}
