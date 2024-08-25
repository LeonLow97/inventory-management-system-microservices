package ports

import (
	"context"

	"github.com/LeonLow97/internal/core/domain"
)

type InventoryRepo interface {
	GetProducts(ctx context.Context, userID int) (*[]domain.Product, error)
	GetProductByID(ctx context.Context, userID, productID int) (*domain.Product, error)
	CreateProduct(ctx context.Context, req domain.Product, userID int) error
	UpdateProduct(ctx context.Context, req domain.Product, userID, productID int) error
	DeleteProduct(ctx context.Context, userID, productID int) error
}
