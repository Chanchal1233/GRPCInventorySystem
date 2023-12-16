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

type ProductService struct {
	productRepo  *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
}

func NewProductService(productRepo *repository.ProductRepository, categoryRepo *repository.CategoryRepository) *ProductService {
	return &ProductService{productRepo: productRepo, categoryRepo: categoryRepo}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	product.ID = uuid.New()
	exists, err := s.productRepo.ExistsBySKUOrName(ctx, product.SKU, product.Name)
	if err != nil {
		log.Printf("Error checking existence: %v", err)
		return nil, status.Errorf(codes.Internal, "error checking existence: %v", err)
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "product with the same SKU or Name already exists")
	}
	categoryExists, err := s.categoryRepo.CategoryExists(ctx, product.CategoryID.String())
	if err != nil {
		log.Printf("Error checking category existence: %v", err)
		return nil, status.Errorf(codes.Internal, "error checking category existence: %v", err)
	}
	if !categoryExists {
		return nil, status.Error(codes.NotFound, "category with the provided ID does not exist")
	}
	err = s.productRepo.CreateProduct(ctx, product)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		return nil, status.Errorf(codes.Internal, "error creating product: %v", err)
	}
	log.Println("Product created successfully:", product.ID)
	return product, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	product, err := s.productRepo.GetProduct(ctx, id.String())
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		log.Printf("Error retrieving product: %v", err)
		return nil, status.Errorf(codes.Internal, "error retrieving product: %v", err)
	}
	return product, nil
}

func (s *ProductService) ListProducts(ctx context.Context) ([]*model.Product, error) {
	products, err := s.productRepo.ListProducts(ctx)
	if err != nil {
		log.Printf("Error listing products: %v", err)
		return nil, status.Errorf(codes.Internal, "error listing products: %v", err)
	}
	log.Println("Products listed successfully")
	return products, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *model.Product) error {
	exists, err := s.productRepo.ExistsBySKUOrNameExcludingID(ctx, product.SKU, product.Name, product.ID)
	if err != nil {
		log.Printf("Error checking for existing product: %v", err)
		return status.Errorf(codes.Internal, "error checking existence: %v", err)
	}
	if exists {
		return status.Error(codes.AlreadyExists, "another product with the same SKU or Name already exists")
	}
	err = s.productRepo.UpdateProduct(ctx, product)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			return status.Error(codes.NotFound, "product not found")
		}
		log.Printf("Error updating product: %v", err)
		return status.Errorf(codes.Internal, "error updating product: %v", err)
	}
	log.Println("Product updated successfully:", product.ID)
	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	err := s.productRepo.DeleteProduct(ctx, id.String())
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			return status.Error(codes.NotFound, "product not found")
		}
		log.Printf("Error deleting product: %v", err)
		return status.Errorf(codes.Internal, "error deleting product: %v", err)
	}
	log.Println("Product deleted successfully:", id)
	return nil
}
