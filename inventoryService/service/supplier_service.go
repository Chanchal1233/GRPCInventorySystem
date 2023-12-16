package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"inventoryService/model"
	"inventoryService/repository"
	"log"
)

type SupplierService struct {
	repo *repository.SupplierRepository
}

func NewSupplierService(repo *repository.SupplierRepository) *SupplierService {
	return &SupplierService{repo: repo}
}

func (s *SupplierService) CreateSupplier(ctx context.Context, supplier *model.Supplier) (*model.Supplier, error) {
	supplier.ID = uuid.New()
	exists, err := s.repo.ExistsByUUID(ctx, supplier.ID)
	if err != nil {
		log.Printf("Error checking existence: %v", err)
		return nil, err
	}
	if exists {
		return nil, errors.New("supplier with the same UUID already exists")
	}
	err = s.repo.CreateSupplier(ctx, supplier)
	if err != nil {
		log.Printf("Error creating supplier: %v", err)
		return nil, err
	}
	log.Println("Supplier created successfully:", supplier.ID)
	return supplier, nil
}

func (s *SupplierService) GetSupplier(ctx context.Context, id uuid.UUID) (*model.Supplier, error) {
	supplier, err := s.repo.GetSupplier(ctx, id.String())
	if err != nil {
		log.Printf("Error retrieving supplier: %v", err)
		return nil, err
	}
	if supplier == nil {
		return nil, errors.New("supplier not found")
	}
	return supplier, nil
}

func (s *SupplierService) ListSuppliers(ctx context.Context) ([]*model.Supplier, error) {
	suppliers, err := s.repo.ListSuppliers(ctx)
	if err != nil {
		log.Printf("Error listing suppliers: %v", err)
		return nil, err
	}
	log.Println("Suppliers listed successfully")
	return suppliers, nil
}

func (s *SupplierService) UpdateSupplier(ctx context.Context, supplier *model.Supplier) error {
	exists, err := s.repo.ExistsByNameExcludingUUID(ctx, supplier.Name, supplier.ID)
	if err != nil {
		log.Printf("Error checking for existing supplier: %v", err)
		return err
	}
	if exists {
		return errors.New("another supplier with the same name already exists")
	}
	err = s.repo.UpdateSupplier(ctx, supplier)
	if err != nil {
		log.Printf("Error updating supplier: %v", err)
		return err
	}
	log.Println("Supplier updated successfully:", supplier.ID)
	return nil
}

func (s *SupplierService) DeleteSupplier(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteSupplier(ctx, id.String())
	if err != nil {
		log.Printf("Error deleting supplier: %v", err)
		return err
	}
	log.Println("Supplier deleted successfully:", id)
	return nil
}
