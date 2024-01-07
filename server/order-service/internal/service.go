package order

import (
	"encoding/json"
	"log"

	"github.com/LeonLow97/internal/kafkago"
	"github.com/segmentio/kafka-go"
)

type Service interface {
	GetOrders(req GetOrdersDTO) (*[]Order, error)
	GetOrderByID(req GetOrderDTO) (*Order, error)
	CreateOrder(req CreateOrderDTO) error
}

type service struct {
	repo              Repository
	segmentioInstance *kafkago.Segmentio
	kafkaConfig       *kafkago.KafkaConfig
}

func NewService(repo Repository, segmentioInstance *kafkago.Segmentio, kafkaConfig *kafkago.KafkaConfig) Service {
	return &service{
		repo:              repo,
		segmentioInstance: segmentioInstance,
		kafkaConfig:       kafkaConfig,
	}
}

func (s service) GetOrders(req GetOrdersDTO) (*[]Order, error) {
	orders, err := s.repo.GetOrders(req)
	if err != nil {
		log.Println("error getting orders", err)
		return nil, err
	}

	return orders, err
}

func (s service) GetOrderByID(req GetOrderDTO) (*Order, error) {
	order, err := s.repo.GetOrderByID(req)
	if err != nil {
		log.Printf("error getting 1 order by order_id %d with error %v\n", req.OrderID, err)
		return nil, err
	}

	return order, nil
}

func (s service) CreateOrder(req CreateOrderDTO) error {
	// get product names and category names via grpc to inventory microservice

	orderEvent := OrderEvent{
		Action:   "create_order",
		UserID:   req.UserID,
		Quantity: req.Quantity,
	}

	jsonData, err := json.Marshal(orderEvent)
	if err != nil {
		return err
	}

	// produce an order to inventory microservice to update inventory count
	createOrderEvent := []kafka.Message{
		{
			Key:   []byte(orderEvent.Action),
			Value: []byte(jsonData),
		},
	}

	go func() {
		if err := s.segmentioInstance.Producer(s.kafkaConfig.BrokerAddress, s.kafkaConfig.TopicName, createOrderEvent); err != nil {
			log.Printf("failed to produce message for %s topic: %v\n", s.kafkaConfig.TopicName, err)
		}
	}()

	// create order with status 'SUBMITTED'

	return nil
}
