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

type InventoryItemService struct {
	repo          *repository.InventoryItemRepository
	productRepo   *repository.ProductRepository
	warehouseRepo *repository.WarehouseRepository
}

func NewInventoryItemService(repo *repository.InventoryItemRepository, productRepo *repository.ProductRepository, warehouseRepo *repository.WarehouseRepository) *InventoryItemService {
	return &InventoryItemService{
		repo:          repo,
		productRepo:   productRepo,
		warehouseRepo: warehouseRepo,
	}
}

func (s *InventoryItemService) CreateInventoryItem(ctx context.Context, item *model.InventoryItem) (*model.InventoryItem, error) {
	item.ID = uuid.New()

	exists, err := s.productRepo.ProductExists(ctx, item.ProductID.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error checking product existence: %v", err)
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "product not found")
	}
	exists, err = s.warehouseRepo.WarehouseExists(ctx, item.WarehouseID.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error checking warehouse existence: %v", err)
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "warehouse not found")
	}

	err = s.repo.CreateInventoryItem(ctx, item)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating inventory item: %v", err)
	}

	return item, nil
}

func (s *InventoryItemService) GetInventoryItem(ctx context.Context, id uuid.UUID) (*model.InventoryItem, error) {
	item, err := s.repo.GetInventoryItem(ctx, id.String())
	if err != nil {
		if errors.Is(err, repository.ErrInventoryItemNotFound) {
			return nil, status.Error(codes.NotFound, "inventory item not found")
		}
		log.Printf("Error retrieving inventory item: %v", err)
		return nil, status.Errorf(codes.Internal, "error retrieving inventory item: %v", err)
	}
	return item, nil
}

func (s *InventoryItemService) ListInventoryItems(ctx context.Context) ([]*model.InventoryItem, error) {
	items, err := s.repo.ListInventoryItems(ctx)
	if err != nil {
		log.Printf("Error listing inventory items: %v", err)
		return nil, status.Errorf(codes.Internal, "error listing inventory items: %v", err)
	}
	log.Println("Inventory items listed successfully")
	return items, nil
}

func (s *InventoryItemService) UpdateInventoryItem(ctx context.Context, item *model.InventoryItem) error {
	exists, err := s.productRepo.ProductExists(ctx, item.ProductID.String())
	if err != nil {
		log.Printf("Error checking product existence: %v", err)
		return status.Errorf(codes.Internal, "error checking product existence: %v", err)
	}
	if !exists {
		return status.Error(codes.NotFound, "product not found")
	}
	exists, err = s.warehouseRepo.WarehouseExists(ctx, item.WarehouseID.String())
	if err != nil {
		log.Printf("Error checking warehouse existence: %v", err)
		return status.Errorf(codes.Internal, "error checking warehouse existence: %v", err)
	}
	if !exists {
		return status.Error(codes.NotFound, "warehouse not found")
	}
	err = s.repo.UpdateInventoryItem(ctx, item)
	if err != nil {
		log.Printf("Error updating inventory item: %v", err)
		return status.Errorf(codes.Internal, "error updating inventory item: %v", err)
	}
	log.Println("Inventory item updated successfully:", item.ID)
	return nil
}

func (s *InventoryItemService) DeleteInventoryItem(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteInventoryItem(ctx, id.String())
	if err != nil {
		if errors.Is(err, repository.ErrInventoryItemNotFound) {
			return status.Error(codes.NotFound, "inventory item not found")
		}
		log.Printf("Error deleting inventory item: %v", err)
		return status.Errorf(codes.Internal, "error deleting inventory item: %v", err)
	}
	log.Println("Inventory item deleted successfully:", id)
	return nil
}
