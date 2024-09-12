package services

import (
	"context"
	"log"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/pkg/config"
	"github.com/LeonLow97/internal/pkg/kafkago"
	"github.com/LeonLow97/internal/ports"
	"github.com/google/uuid"
)

type Service interface {
	GetOrders(userID int) (*[]domain.Order, error)
	GetOrderByID(userID, orderID int) (*domain.Order, error)
	CreateOrder(ctx context.Context, req domain.Order, userID int, productName string) error
}

type service struct {
	cfg  config.Config
	repo ports.Repository
}

func NewService(cfg config.Config, r ports.Repository) Service {
	return &service{
		cfg:  cfg,
		repo: r,
	}
}

func (s service) GetOrders(userID int) (*[]domain.Order, error) {
	orders, err := s.repo.GetOrders(userID)
	if err != nil {
		log.Printf("failed to get orders for user id %d with error: %v\n", userID, err)
		return nil, err
	}

	return orders, nil
}

func (s service) GetOrderByID(userID, orderID int) (*domain.Order, error) {
	order, err := s.repo.GetOrderByID(userID, orderID)
	if err != nil {
		log.Printf("failed to get 1 order by order id %d for user id %d with error %v\n", orderID, userID, err)
		return nil, err
	}

	return order, nil
}

func (s service) CreateOrder(ctx context.Context, req domain.Order, userID int, productName string) error {
	// retrieve product id
	productID, err := s.repo.GetProductID(ctx, userID, req.BrandName, req.CategoryName, productName)
	if err != nil {
		log.Printf("failed to get product id with error: %v\n", err)
		return err
	}

	// generate uuid for order event
	orderUUID := uuid.New().String()

	producerErrorChan := make(chan error)
	go func() {
		producerErrorChan <- s.repo.ProduceOrderMessage(s.cfg.KafkaConfig.BrokerAddress, kafkago.TOPIC_DECREMENT_INVENTORY, orderUUID, productID, userID, req.Quantity)
		close(producerErrorChan)
	}()

	// create order with status 'SUBMITTED'
	req.Status = "SUBMITTED"
	req.OrderUUID = orderUUID
	if err := s.repo.CreateOrder(req, userID, productID); err != nil {
		log.Printf("failed to create order with error: %v\n", err)
		return err
	}

	// wait for the goroutine to finish and capture the error (if any)
	if err := <-producerErrorChan; err != nil {
		log.Printf("failed to produce order message with error: %v\n", err)
		return err
	}

	return nil
}
