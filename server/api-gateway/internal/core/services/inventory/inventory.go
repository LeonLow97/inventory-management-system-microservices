package inventory

import (
	"context"
	"log"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/ports"
)

type Inventory interface {
	GetProducts(ctx context.Context, userID int) (*[]domain.Product, error)
	GetProductByID(ctx context.Context, userID, productID int) (*domain.Product, error)
	CreateProduct(ctx context.Context, req domain.Product, userID int) error
	UpdateProduct(ctx context.Context, req domain.Product, userID, productID int) error
	DeleteProduct(ctx context.Context, userID, productID int) error
}

type service struct {
	inventoryRepo ports.InventoryRepo
}

func NewInventoryService(r ports.InventoryRepo) Inventory {
	return &service{
		inventoryRepo: r,
	}
}

func (s *service) GetProducts(ctx context.Context, userID int) (*[]domain.Product, error) {
	products, err := s.inventoryRepo.GetProducts(ctx, userID)
	if err != nil {
		log.Printf("failed to get products with error: %v\n", err)
		return nil, err
	}

	return products, nil
}

func (s *service) GetProductByID(ctx context.Context, userID, productID int) (*domain.Product, error) {
	product, err := s.inventoryRepo.GetProductByID(ctx, userID, productID)
	if err != nil {
		log.Printf("failed to get product by ID with error: %v\n", err)
		return nil, err
	}

	return product, nil
}

func (s *service) CreateProduct(ctx context.Context, req domain.Product, userID int) error {
	return s.inventoryRepo.CreateProduct(ctx, req, userID)
}

func (s *service) UpdateProduct(ctx context.Context, req domain.Product, userID, productID int) error {
	return s.inventoryRepo.UpdateProduct(ctx, req, userID, productID)
}

func (s *service) DeleteProduct(ctx context.Context, userID, productID int) error {
	return s.inventoryRepo.DeleteProduct(ctx, userID, productID)
}
