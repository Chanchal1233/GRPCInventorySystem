package model

import (
	"github.com/google/uuid"
	"time"
)

type StockMovement struct {
	ID                     uuid.UUID         `json:"id"`
	InventoryItemID        uuid.UUID         `json:"inventory_item_id"`
	Type                   StockMovementType `json:"type"`
	Quantity               int               `json:"quantity"`
	Date                   time.Time         `json:"date"`
	SourceWarehouseID      uuid.UUID         `json:"source_warehouse_id"`
	DestinationWarehouseID uuid.UUID         `json:"destination_warehouse_id"`
}

type StockMovementType int

const (
	Addition StockMovementType = iota
	Removal
	Transfer
)
