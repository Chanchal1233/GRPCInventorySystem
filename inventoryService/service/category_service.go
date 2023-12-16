package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"inventoryService/model"
	"inventoryService/repository"
	"log"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error) {
	category.ID = uuid.New()
	exists, err := s.repo.ExistsByName(ctx, category.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error checking existence: %v", err)
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "category with the same name already exists")
	}
	err = s.repo.CreateCategory(ctx, category)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating category: %v", err)
	}
	return category, nil
}

func (s *CategoryService) GetCategory(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	category, err := s.repo.GetCategory(ctx, id.String())
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "category not found")
		}
		log.Printf("Error retrieving category: %v", err)
		return nil, status.Errorf(codes.Internal, "error retrieving category: %v", err)
	}
	return category, nil
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]*model.Category, error) {
	categories, err := s.repo.ListCategories(ctx)
	if err != nil {
		log.Printf("Error listing categories: %v", err)
		return nil, status.Errorf(codes.Internal, "error listing categories: %v", err)
	}
	log.Println("Categories listed successfully")
	return categories, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, category *model.Category) error {
	exists, err := s.repo.ExistsByNameExcludingID(ctx, category.Name, category.ID)
	if err != nil {
		log.Printf("Error checking for existing category: %v", err)
		return status.Errorf(codes.Internal, "error checking existence: %v", err)
	}
	if exists {
		return status.Error(codes.AlreadyExists, "another category with the same Name already exists")
	}
	err = s.repo.UpdateCategory(ctx, category)
	if err != nil {
		log.Printf("Error updating category: %v", err)
		return status.Errorf(codes.Internal, "error updating category: %v", err)
	}
	log.Println("Category updated successfully:", category.ID)
	return nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteCategory(ctx, id.String())
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			return status.Error(codes.NotFound, "category not found")
		}
		log.Printf("Error deleting category: %v", err)
		return status.Errorf(codes.Internal, "error deleting category: %v", err)
	}
	log.Println("Category deleted successfully:", id)
	return nil
}
