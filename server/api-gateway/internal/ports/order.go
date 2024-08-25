package ports

import (
	"context"

	"github.com/LeonLow97/internal/core/domain"
)

type OrderRepo interface {
	GetOrders(ctx context.Context, userID int) (*[]domain.Order, error)
	GetOrder(ctx context.Context, userID, orderID int) (*domain.Order, error)
	CreateOrder(ctx context.Context, req domain.Order, userID int) error
}
