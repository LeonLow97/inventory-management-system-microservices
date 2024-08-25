package order

import (
	"context"
	"log"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/ports"
)

type Order interface {
	GetOrders(ctx context.Context, userID int) (*[]domain.Order, error)
	GetOrder(ctx context.Context, userID, orderID int) (*domain.Order, error)
	CreateOrder(ctx context.Context, req domain.Order, userID int) error
}

type service struct {
	orderRepo ports.OrderRepo
}

func NewOrderService(r ports.OrderRepo) Order {
	return &service{
		orderRepo: r,
	}
}

func (s *service) GetOrders(ctx context.Context, userID int) (*[]domain.Order, error) {
	orders, err := s.orderRepo.GetOrders(ctx, userID)
	if err != nil {
		log.Printf("failed to get orders with error: %v\n", err)
		return nil, err
	}

	return orders, nil
}

func (s *service) GetOrder(ctx context.Context, userID, orderID int) (*domain.Order, error) {
	order, err := s.orderRepo.GetOrder(ctx, userID, orderID)
	if err != nil {
		log.Printf("failed to get order with error: %v\n", err)
		return nil, err
	}

	return order, nil
}

func (s *service) CreateOrder(ctx context.Context, req domain.Order, userID int) error {
	return s.orderRepo.CreateOrder(ctx, req, userID)
}
