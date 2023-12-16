package repository

import (
	"context"
	"errors"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"inventoryService/model"
	"log"
)

type ProductRepository struct {
	session *gocql.Session
}

var ErrProductNotFound = errors.New("product not found")

func NewProductRepository(session *gocql.Session) *ProductRepository {
	return &ProductRepository{session: session}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *model.Product) error {
	return r.session.Query(`INSERT INTO products (id, name, description, category_id, price, sku) VALUES (?, ?, ?, ?, ?, ?)`,
		product.ID.String(), product.Name, product.Description, product.CategoryID.String(), product.Price, product.SKU).WithContext(ctx).Exec()
}

func (r *ProductRepository) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	product := &model.Product{}
	var productID, categoryID string
	if err := r.session.Query(`SELECT id, name, description, category_id, price, sku FROM products WHERE id = ? LIMIT 1`,
		id).WithContext(ctx).Consistency(gocql.One).Scan(&productID, &product.Name, &product.Description, &categoryID, &product.Price, &product.SKU); err != nil {
		return nil, err
	}
	parsedUUID, err := uuid.Parse(productID)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		return nil, err
	}
	product.ID = parsedUUID
	catUUID, err := uuid.Parse(categoryID)
	if err != nil {
		log.Printf("Error parsing category UUID: %v", err)
		return nil, err
	}
	product.CategoryID = catUUID

	return product, nil
}

func (r *ProductRepository) ListProducts(ctx context.Context) ([]*model.Product, error) {
	var products []*model.Product
	iter := r.session.Query(`SELECT id, name, description, category_id, price, sku FROM products`).WithContext(ctx).Iter()
	var id, categoryID, name, description, sku string
	var price float64
	for iter.Scan(&id, &name, &description, &categoryID, &price, &sku) {
		convertedID, err := uuid.Parse(id)
		if err != nil {
			continue
		}
		convertedCategoryID, err := uuid.Parse(categoryID)
		if err != nil {
			continue
		}
		product := &model.Product{
			ID:          convertedID,
			Name:        name,
			Description: description,
			CategoryID:  convertedCategoryID,
			Price:       price,
			SKU:         sku,
		}
		products = append(products, product)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product *model.Product) error {
	return r.session.Query(`UPDATE products SET name = ?, description = ?, category_id = ?, price = ?, sku = ? WHERE id = ?`,
		product.Name, product.Description, product.CategoryID.String(), product.Price, product.SKU, product.ID.String()).WithContext(ctx).Exec()
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	return r.session.Query(`DELETE FROM products WHERE id = ?`, id).WithContext(ctx).Exec()
}

func (r *ProductRepository) ExistsBySKUOrName(ctx context.Context, sku, name string) (bool, error) {
	var count int
	if err := r.session.Query(`SELECT COUNT(*) FROM products WHERE sku = ?`, sku).WithContext(ctx).Consistency(gocql.One).Scan(&count); err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	if err := r.session.Query(`SELECT COUNT(*) FROM products WHERE name = ?`, name).WithContext(ctx).Consistency(gocql.One).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ProductRepository) ExistsBySKUOrNameExcludingID(ctx context.Context, sku, name string, id uuid.UUID) (bool, error) {
	iter := r.session.Query(`SELECT id FROM products WHERE sku = ?`, sku).WithContext(ctx).Consistency(gocql.One).Iter()
	var productID string
	for iter.Scan(&productID) {
		if productID != id.String() {
			return true, nil
		}
	}
	if err := iter.Close(); err != nil {
		return false, err
	}
	iter = r.session.Query(`SELECT id FROM products WHERE name = ?`, name).WithContext(ctx).Consistency(gocql.One).Iter()
	for iter.Scan(&productID) {
		if productID != id.String() {
			return true, nil
		}
	}
	if err := iter.Close(); err != nil {
		return false, err
	}
	return false, nil
}

func (r *ProductRepository) ProductExists(ctx context.Context, id string) (bool, error) {
	var count int
	if err := r.session.Query(`SELECT COUNT(*) FROM products WHERE id = ?`, id).WithContext(ctx).Consistency(gocql.One).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
