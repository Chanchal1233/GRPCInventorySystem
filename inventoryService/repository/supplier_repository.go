package repository

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"inventoryService/model"
)

type SupplierRepository struct {
	session *gocql.Session
}

func NewSupplierRepository(session *gocql.Session) *SupplierRepository {
	return &SupplierRepository{session: session}
}

func (r *SupplierRepository) CreateSupplier(ctx context.Context, supplier *model.Supplier) error {
	return r.session.Query(`INSERT INTO suppliers (id, name, contact_info) VALUES (?, ?, ?)`,
		supplier.ID.String(), supplier.Name, supplier.ContactInfo).WithContext(ctx).Exec()
}

func (r *SupplierRepository) GetSupplier(ctx context.Context, id string) (*model.Supplier, error) {
	var idStr string
	supplier := &model.Supplier{}
	if err := r.session.Query(`SELECT id, name, contact_info FROM suppliers WHERE id = ? LIMIT 1`,
		id).WithContext(ctx).Consistency(gocql.One).Scan(&idStr, &supplier.Name, &supplier.ContactInfo); err != nil {
		return nil, err
	}
	supplier.ID, _ = uuid.Parse(idStr)
	return supplier, nil
}

func (r *SupplierRepository) ListSuppliers(ctx context.Context) ([]*model.Supplier, error) {
	var suppliers []*model.Supplier
	iter := r.session.Query(`SELECT id, name, contact_info FROM suppliers`).WithContext(ctx).Iter()
	var idStr, name, contactInfo string
	for iter.Scan(&idStr, &name, &contactInfo) {
		supplier := &model.Supplier{
			ID:          uuid.MustParse(idStr),
			Name:        name,
			ContactInfo: contactInfo,
		}
		suppliers = append(suppliers, supplier)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return suppliers, nil
}

func (r *SupplierRepository) UpdateSupplier(ctx context.Context, supplier *model.Supplier) error {
	return r.session.Query(`UPDATE suppliers SET name = ?, contact_info = ? WHERE id = ?`,
		supplier.Name, supplier.ContactInfo, supplier.ID.String()).WithContext(ctx).Exec()
}

func (r *SupplierRepository) DeleteSupplier(ctx context.Context, id string) error {
	return r.session.Query(`DELETE FROM suppliers WHERE id = ?`, id).WithContext(ctx).Exec()
}

func (r *SupplierRepository) ExistsByUUID(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int
	if err := r.session.Query(`SELECT COUNT(*) FROM suppliers WHERE id = ?`, id.String()).WithContext(ctx).Consistency(gocql.One).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *SupplierRepository) ExistsByNameExcludingUUID(ctx context.Context, name string, id uuid.UUID) (bool, error) {
	iter := r.session.Query(`SELECT id FROM suppliers WHERE name = ?`, name).WithContext(ctx).Consistency(gocql.One).Iter()
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
