package model

import "github.com/google/uuid"

type InventoryItem struct {
	ID              uuid.UUID `json:"id"`
	ProductID       uuid.UUID `json:"product_id"`
	WarehouseID     uuid.UUID `json:"warehouse_id"`
	Quantity        int       `json:"quantity"`
	ReorderLevel    int       `json:"reorder_level"`
	ReorderQuantity int       `json:"reorder_quantity"`
}
