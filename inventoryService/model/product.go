package model

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CategoryID  uuid.UUID `json:"category_id"`
	Price       float64   `json:"price"`
	SKU         string    `json:"sku"`
}
