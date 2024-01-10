package main

import (
	"net/http"

	order "github.com/LeonLow97/internal"
	grpcclient "github.com/LeonLow97/internal/grpc"
	kafkago "github.com/LeonLow97/pkg/kafkago"
	s3client "github.com/LeonLow97/pkg/s3"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

const (
	INVENTORY_SERVICE_URL = "inventory-service:8002"
)

func (app *application) setupDBDependencies(db *sqlx.DB, segmentioInstance *kafkago.Segmentio, clients *grpcClientConn, kafkaConfigUpdateInventoryCount *kafkago.KafkaConfig, s3Session s3client.BucketClient) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	inventoryServiceClient := grpcclient.NewInventoryGRPCClient(clients.inventoryConn)

	orderRepo := order.NewRepository(db)
	segmentioInstance.AddTopicConfig(kafkaConfigUpdateInventoryCount.TopicName, 1, 1)
	orderService := order.NewService(orderRepo, segmentioInstance, inventoryServiceClient, kafkaConfigUpdateInventoryCount)
	orderServiceServer := order.NewOrderGRPCHandler(orderService)
	app.orderService = orderServiceServer

	return r
}
