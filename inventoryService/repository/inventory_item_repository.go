package repository

import (
	"context"
	"errors"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"inventoryService/model"
)

type InventoryItemRepository struct {
	session *gocql.Session
}

func NewInventoryItemRepository(session *gocql.Session) *InventoryItemRepository {
	return &InventoryItemRepository{session: session}
}

func (r *InventoryItemRepository) CreateInventoryItem(ctx context.Context, item *model.InventoryItem) error {
	return r.session.Query(`INSERT INTO inventory_items (id, product_id, warehouse_id, quantity, reorder_level, reorder_quantity) VALUES (?, ?, ?, ?, ?, ?)`,
		item.ID.String(), item.ProductID.String(), item.WarehouseID.String(), item.Quantity, item.ReorderLevel, item.ReorderQuantity).WithContext(ctx).Exec()
}

var ErrInventoryItemNotFound = errors.New("inventory item not found")

func (r *InventoryItemRepository) GetInventoryItem(ctx context.Context, id string) (*model.InventoryItem, error) {
	var idStr, productIdStr, warehouseIdStr string
	item := &model.InventoryItem{}
	err := r.session.Query(`SELECT id, product_id, warehouse_id, quantity, reorder_level, reorder_quantity FROM inventory_items WHERE id = ? LIMIT 1`,
		id).WithContext(ctx).Consistency(gocql.One).Scan(&idStr, &productIdStr, &warehouseIdStr, &item.Quantity, &item.ReorderLevel, &item.ReorderQuantity)
	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, ErrInventoryItemNotFound
		}
		return nil, err
	}
	item.ID, _ = uuid.Parse(idStr)
	item.ProductID, _ = uuid.Parse(productIdStr)
	item.WarehouseID, _ = uuid.Parse(warehouseIdStr)
	return item, nil
}

func (r *InventoryItemRepository) ListInventoryItems(ctx context.Context) ([]*model.InventoryItem, error) {
	var items []*model.InventoryItem
	iter := r.session.Query(`SELECT id, product_id, warehouse_id, quantity, reorder_level, reorder_quantity FROM inventory_items`).WithContext(ctx).Iter()
	var item model.InventoryItem
	for iter.Scan(&item.ID, &item.ProductID, &item.WarehouseID, &item.Quantity, &item.ReorderLevel, &item.ReorderQuantity) {
		items = append(items, &item)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *InventoryItemRepository) UpdateInventoryItem(ctx context.Context, item *model.InventoryItem) error {
	return r.session.Query(`UPDATE inventory_items SET product_id = ?, warehouse_id = ?, quantity = ?, reorder_level = ?, reorder_quantity = ? WHERE id = ?`,
		item.ProductID.String(), item.WarehouseID.String(), item.Quantity, item.ReorderLevel, item.ReorderQuantity, item.ID.String()).WithContext(ctx).Exec()
}

func (r *InventoryItemRepository) DeleteInventoryItem(ctx context.Context, id string) error {
	return r.session.Query(`DELETE FROM inventory_items WHERE id = ?`, id).WithContext(ctx).Exec()
}
