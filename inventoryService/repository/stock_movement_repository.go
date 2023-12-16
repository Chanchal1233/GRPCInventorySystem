package repository

import (
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"inventoryService/model"
	"time"
)

type StockMovementRepository struct {
	session *gocql.Session
}

func NewStockMovementRepository(session *gocql.Session) *StockMovementRepository {
	return &StockMovementRepository{session: session}
}

func (r *StockMovementRepository) CreateStockMovement(ctx context.Context, movement *model.StockMovement) error {
	return r.session.Query(`INSERT INTO stock_movements (id, inventory_item_id, type, quantity, date, source_warehouse_id, destination_warehouse_id) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		movement.ID.String(), movement.InventoryItemID.String(), movement.Type, movement.Quantity, movement.Date, movement.SourceWarehouseID.String(), movement.DestinationWarehouseID.String()).WithContext(ctx).Exec()
}

func (r *StockMovementRepository) GetStockMovement(ctx context.Context, id string) (*model.StockMovement, error) {
	var idStr, inventoryItemIdStr, sourceWarehouseIdStr, destinationWarehouseIdStr string
	movement := &model.StockMovement{}
	if err := r.session.Query(`SELECT id, inventory_item_id, type, quantity, date, source_warehouse_id, destination_warehouse_id FROM stock_movements WHERE id = ? LIMIT 1`,
		id).WithContext(ctx).Consistency(gocql.One).Scan(&idStr, &inventoryItemIdStr, &movement.Type, &movement.Quantity, &movement.Date, &sourceWarehouseIdStr, &destinationWarehouseIdStr); err != nil {
		return nil, err
	}
	movement.ID, _ = uuid.Parse(idStr)
	movement.InventoryItemID, _ = uuid.Parse(inventoryItemIdStr)
	movement.SourceWarehouseID, _ = uuid.Parse(sourceWarehouseIdStr)
	movement.DestinationWarehouseID, _ = uuid.Parse(destinationWarehouseIdStr)
	return movement, nil
}

func (r *StockMovementRepository) ListStockMovements(ctx context.Context) ([]*model.StockMovement, error) {
	var movements []*model.StockMovement
	iter := r.session.Query(`SELECT id, inventory_item_id, type, quantity, date, source_warehouse_id, destination_warehouse_id FROM stock_movements`).WithContext(ctx).Iter()
	var idStr, inventoryItemIdStr, sourceWarehouseIdStr, destinationWarehouseIdStr string
	var movementTypeInt int
	var quantity int
	var date time.Time
	for iter.Scan(&idStr, &inventoryItemIdStr, &movementTypeInt, &quantity, &date, &sourceWarehouseIdStr, &destinationWarehouseIdStr) {
		if movementTypeInt < 0 || movementTypeInt > 2 {
			return nil, fmt.Errorf("invalid stock movement type: %d", movementTypeInt)
		}
		movement := &model.StockMovement{
			ID:                     uuid.MustParse(idStr),
			InventoryItemID:        uuid.MustParse(inventoryItemIdStr),
			Type:                   model.StockMovementType(movementTypeInt),
			Quantity:               quantity,
			Date:                   date,
			SourceWarehouseID:      uuid.MustParse(sourceWarehouseIdStr),
			DestinationWarehouseID: uuid.MustParse(destinationWarehouseIdStr),
		}
		movements = append(movements, movement)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return movements, nil
}

func (r *StockMovementRepository) UpdateStockMovement(ctx context.Context, movement *model.StockMovement) error {
	return r.session.Query(`UPDATE stock_movements SET inventory_item_id = ?, type = ?, quantity = ?, date = ?, source_warehouse_id = ?, destination_warehouse_id = ? WHERE id = ?`,
		movement.InventoryItemID.String(), movement.Type, movement.Quantity, movement.Date, movement.SourceWarehouseID.String(), movement.DestinationWarehouseID.String(), movement.ID.String()).WithContext(ctx).Exec()
}

func (r *StockMovementRepository) DeleteStockMovement(ctx context.Context, id string) error {
	return r.session.Query(`DELETE FROM stock_movements WHERE id = ?`, id).WithContext(ctx).Exec()
}
