package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"inventoryService/model"
	"inventoryService/repository"
	"log"
)

type StockMovementService struct {
	repo *repository.StockMovementRepository
}

func NewStockMovementService(repo *repository.StockMovementRepository) *StockMovementService {
	return &StockMovementService{repo: repo}
}

func (s *StockMovementService) CreateStockMovement(ctx context.Context, movement *model.StockMovement) (*model.StockMovement, error) {
	movement.ID = uuid.New()
	err := s.repo.CreateStockMovement(ctx, movement)
	if err != nil {
		log.Printf("Error creating stock movement: %v", err)
		return nil, err
	}
	log.Println("Stock movement created successfully:", movement.ID)
	return movement, nil
}

func (s *StockMovementService) GetStockMovement(ctx context.Context, id uuid.UUID) (*model.StockMovement, error) {
	movement, err := s.repo.GetStockMovement(ctx, id.String())
	if err != nil {
		log.Printf("Error retrieving stock movement: %v", err)
		return nil, err
	}
	if movement == nil {
		return nil, errors.New("stock movement not found")
	}
	return movement, nil
}

func (s *StockMovementService) ListStockMovements(ctx context.Context) ([]*model.StockMovement, error) {
	movements, err := s.repo.ListStockMovements(ctx)
	if err != nil {
		log.Printf("Error listing stock movements: %v", err)
		return nil, err
	}
	log.Println("Stock movements listed successfully")
	return movements, nil
}

func (s *StockMovementService) UpdateStockMovement(ctx context.Context, movement *model.StockMovement) error {
	err := s.repo.UpdateStockMovement(ctx, movement)
	if err != nil {
		log.Printf("Error updating stock movement: %v", err)
		return err
	}
	log.Println("Stock movement updated successfully:", movement.ID)
	return nil
}

func (s *StockMovementService) DeleteStockMovement(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteStockMovement(ctx, id.String())
	if err != nil {
		log.Printf("Error deleting stock movement: %v", err)
		return err
	}
	log.Println("Stock movement deleted successfully:", id)
	return nil
}
