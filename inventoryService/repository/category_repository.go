package repository

import (
	"context"
	"errors"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"inventoryService/model"
	"log"
)

type CategoryRepository struct {
	session *gocql.Session
}

func NewCategoryRepository(session *gocql.Session) *CategoryRepository {
	return &CategoryRepository{session: session}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, category *model.Category) error {
	return r.session.Query(`INSERT INTO categories (id, name, description) VALUES (?, ?, ?)`,
		category.ID.String(), category.Name, category.Description).WithContext(ctx).Exec()
}

var ErrCategoryNotFound = errors.New("category not found")

func (r *CategoryRepository) GetCategory(ctx context.Context, id string) (*model.Category, error) {
	category := &model.Category{}
	var categoryID string
	err := r.session.Query(`SELECT id, name, description FROM categories WHERE id = ? LIMIT 1`,
		id).WithContext(ctx).Consistency(gocql.One).Scan(&categoryID, &category.Name, &category.Description)
	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	parsedUUID, err := uuid.Parse(categoryID)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		return nil, err
	}
	category.ID = parsedUUID
	return category, nil
}

func (r *CategoryRepository) ListCategories(ctx context.Context) ([]*model.Category, error) {
	var categories []*model.Category
	iter := r.session.Query(`SELECT id, name, description FROM categories`).WithContext(ctx).Iter()
	var category model.Category
	for iter.Scan(&category.ID, &category.Name, &category.Description) {
		categories = append(categories, &category)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, category *model.Category) error {
	return r.session.Query(`UPDATE categories SET name = ?, description = ? WHERE id = ?`,
		category.Name, category.Description, category.ID.String()).WithContext(ctx).Exec()
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id string) error {
	exists, err := r.CategoryExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrCategoryNotFound
	}
	err = r.session.Query(`DELETE FROM categories WHERE id = ?`, id).WithContext(ctx).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (r *CategoryRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int
	if err := r.session.Query(`SELECT COUNT(*) FROM categories WHERE name = ?`, name).WithContext(ctx).Consistency(gocql.One).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

var ErrIterClose = errors.New("error closing iterator")

func (r *CategoryRepository) ExistsByNameExcludingID(ctx context.Context, name string, id uuid.UUID) (bool, error) {
	iter := r.session.Query(`SELECT id FROM categories WHERE name = ?`, name).WithContext(ctx).Consistency(gocql.One).Iter()
	var categoryID string
	for iter.Scan(&categoryID) {
		if categoryID != id.String() {
			return true, nil
		}
	}
	if err := iter.Close(); err != nil {
		log.Printf("Error closing iterator: %v", err)
		return false, ErrIterClose
	}
	return false, nil
}

func (r *CategoryRepository) CategoryExists(ctx context.Context, id string) (bool, error) {
	var count int
	if err := r.session.Query(`SELECT COUNT(*) FROM categories WHERE id = ?`, id).WithContext(ctx).Consistency(gocql.One).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
