package repository

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"inventoryService/model"
)

type WarehouseRepository struct {
	session *gocql.Session
}

func NewWarehouseRepository(session *gocql.Session) *WarehouseRepository {
	return &WarehouseRepository{session: session}
}

func (r *WarehouseRepository) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	return r.session.Query(`INSERT INTO warehouses (id, name, location) VALUES (?, ?, ?)`,
		warehouse.ID.String(), warehouse.Name, warehouse.Location).WithContext(ctx).Exec()
}

func (r *WarehouseRepository) GetWarehouse(ctx context.Context, id string) (*model.Warehouse, error) {
	var idStr string
	warehouse := &model.Warehouse{}
	if err := r.session.Query(`SELECT id, name, location FROM warehouses WHERE id = ? LIMIT 1`,
		id).WithContext(ctx).Consistency(gocql.One).Scan(&idStr, &warehouse.Name, &warehouse.Location); err != nil {
		return nil, err
	}
	warehouse.ID, _ = uuid.Parse(idStr)
	return warehouse, nil
}

func (r *WarehouseRepository) ListWarehouses(ctx context.Context) ([]*model.Warehouse, error) {
	var warehouses []*model.Warehouse
	iter := r.session.Query(`SELECT id, name, location FROM warehouses`).WithContext(ctx).Iter()
	var id string
	var name, location string
	for iter.Scan(&id, &name, &location) {
		convertedID, err := uuid.Parse(id)
		if err != nil {
			continue
		}
		warehouse := &model.Warehouse{
			ID:       convertedID,
			Name:     name,
			Location: location,
		}
		warehouses = append(warehouses, warehouse)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return warehouses, nil
}

func (r *WarehouseRepository) UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	return r.session.Query(`UPDATE warehouses SET name = ?, location = ? WHERE id = ?`,
		warehouse.Name, warehouse.Location, warehouse.ID.String()).WithContext(ctx).Exec()
}

func (r *WarehouseRepository) DeleteWarehouse(ctx context.Context, id string) error {
	return r.session.Query(`DELETE FROM warehouses WHERE id = ?`, id).WithContext(ctx).Exec()
}

func (r *WarehouseRepository) ExistsByUUID(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int
	if err := r.session.Query(`SELECT COUNT(*) FROM warehouses WHERE id = ?`, id.String()).WithContext(ctx).Consistency(gocql.One).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *WarehouseRepository) ExistsByNameExcludingUUID(ctx context.Context, name string, id uuid.UUID) (bool, error) {
	iter := r.session.Query(`SELECT id FROM warehouses WHERE name = ?`, name).WithContext(ctx).Consistency(gocql.One).Iter()
	var idStr string
	for iter.Scan(&idStr) {
		if idStr != id.String() {
			return true, nil
		}
	}
	if err := iter.Close(); err != nil {
		return false, err
	}
	return false, nil
}

func (r *WarehouseRepository) WarehouseExists(ctx context.Context, id string) (bool, error) {
	var count int
	if err := r.session.Query(`SELECT COUNT(*) FROM warehouses WHERE id = ?`, id).WithContext(ctx).Consistency(gocql.One).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
