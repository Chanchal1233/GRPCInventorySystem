package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"inventoryService/model"
	"inventoryService/repository"
	"log"
)

type WarehouseService struct {
	repo *repository.WarehouseRepository
}

func NewWarehouseService(repo *repository.WarehouseRepository) *WarehouseService {
	return &WarehouseService{repo: repo}
}

func (s *WarehouseService) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) (*model.Warehouse, error) {
	warehouse.ID = uuid.New()
	exists, err := s.repo.ExistsByUUID(ctx, warehouse.ID)
	if err != nil {
		log.Printf("Error checking existence: %v", err)
		return nil, err
	}
	if exists {
		return nil, errors.New("warehouse with the same UUID already exists")
	}
	err = s.repo.CreateWarehouse(ctx, warehouse)
	if err != nil {
		log.Printf("Error creating warehouse: %v", err)
		return nil, err
	}
	log.Println("Warehouse created successfully:", warehouse.ID)
	return warehouse, nil
}

func (s *WarehouseService) GetWarehouse(ctx context.Context, id uuid.UUID) (*model.Warehouse, error) {
	warehouse, err := s.repo.GetWarehouse(ctx, id.String())
	if err != nil {
		log.Printf("Error retrieving warehouse: %v", err)
		return nil, err
	}
	if warehouse == nil {
		return nil, errors.New("warehouse not found")
	}
	return warehouse, nil
}

func (s *WarehouseService) ListWarehouses(ctx context.Context) ([]*model.Warehouse, error) {
	warehouses, err := s.repo.ListWarehouses(ctx)
	if err != nil {
		log.Printf("Error listing warehouses: %v", err)
		return nil, err
	}
	log.Println("Warehouses listed successfully")
	return warehouses, nil
}

func (s *WarehouseService) UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	exists, err := s.repo.ExistsByNameExcludingUUID(ctx, warehouse.Name, warehouse.ID)
	if err != nil {
		log.Printf("Error checking for existing warehouse: %v", err)
		return err
	}
	if exists {
		return errors.New("another warehouse with the same name already exists")
	}
	err = s.repo.UpdateWarehouse(ctx, warehouse)
	if err != nil {
		log.Printf("Error updating warehouse: %v", err)
		return err
	}
	log.Println("Warehouse updated successfully:", warehouse.ID)
	return nil
}

func (s *WarehouseService) DeleteWarehouse(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteWarehouse(ctx, id.String())
	if err != nil {
		log.Printf("Error deleting warehouse: %v", err)
		return err
	}
	log.Println("Warehouse deleted successfully:", id)
	return nil
}
