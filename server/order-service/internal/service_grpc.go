package order

import (
	"encoding/json"
	"log"

	grpcclient "github.com/LeonLow97/internal/grpc"
	kafkago "github.com/LeonLow97/pkg/kafkago"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type service struct {
	repo              Repository
	segmentioInstance *kafkago.Segmentio
	kafkaConfig       *kafkago.KafkaConfig
	grpcClient        grpcclient.InventoryServiceClient
}

func NewService(repo Repository, segmentioInstance *kafkago.Segmentio, grpcClient grpcclient.InventoryServiceClient, kafkaConfig *kafkago.KafkaConfig) Service {
	return &service{
		repo:              repo,
		segmentioInstance: segmentioInstance,
		kafkaConfig:       kafkaConfig,
		grpcClient:        grpcClient,
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
	// get product names (product_id) and category names via grpc to inventory microservice
	resp, err := s.grpcClient.GRPCGetProductDetailsHandler(req.UserID, req.BrandName, req.CategoryName, req.ProductName)
	if err != nil {
		log.Println(err)
		return err
	}
	if resp.ProductID == 0 {
		return ErrProductNotFound
	}

	// generate uuid for order event
	orderUUID := uuid.New().String()

	orderEvent := kafkago.OrderEvent{
		OrderUUID: orderUUID,
		ProductID: int(resp.ProductID),
		UserID:    req.UserID,
		Quantity:  req.Quantity,
	}

	jsonData, err := json.Marshal(orderEvent)
	if err != nil {
		return err
	}

	// produce an order to inventory microservice to update inventory count
	createOrderEvent := []kafka.Message{
		{
			Key:   []byte(orderEvent.OrderUUID),
			Value: []byte(jsonData),
		},
	}

	go func() {
		if err := s.segmentioInstance.Producer(s.kafkaConfig.BrokerAddress, s.kafkaConfig.TopicName, createOrderEvent); err != nil {
			log.Printf("failed to produce message for %s topic: %v\n", s.kafkaConfig.TopicName, err)
		}
	}()

	// create order with status 'SUBMITTED'
	req.Status = "SUBMITTED"
	req.OrderUUID = orderUUID
	if err := s.repo.CreateOrder(req, int(resp.ProductID)); err != nil {
		return err
	}

	return nil
}
