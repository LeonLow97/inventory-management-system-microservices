package ports

import (
	"context"

	"github.com/LeonLow97/internal/core/domain"
)

type Repository interface {
	// GRPC server (inventory)
	GetProductID(ctx context.Context, userID int, brandName, categoryName, productName string) (int, error)

	// kafka
	ProduceOrderMessage(brokerAddress, topic string, orderID int, productID, userID, orderQuantity int) error
	ConsumeOrderStatus(brokerAddress, topic string)

	// postgres
	GetOrders(userID int) (*[]domain.Order, error)
	GetOrderByID(userID, orderID int) (*domain.Order, error)
	CreateOrder(req domain.Order, userID, productID int) (int, error)
	UpdateOrderByID(req domain.Order) error
}
