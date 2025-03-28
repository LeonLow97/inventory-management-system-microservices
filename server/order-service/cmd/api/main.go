package main

import (
	"log"

	"github.com/LeonLow97/internal/adapters/outbound"
	"github.com/LeonLow97/internal/core/services"
	"github.com/LeonLow97/internal/pkg/config"
	"github.com/LeonLow97/internal/pkg/consul"
	"github.com/LeonLow97/internal/pkg/grpcclient"
	"github.com/LeonLow97/internal/pkg/grpcserver"
	kafkago "github.com/LeonLow97/internal/pkg/kafkago"
	postgres_conn "github.com/LeonLow97/internal/pkg/postgres"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load config", err)
	}

	// initiate kafka-go segmentio instance
	segmentioInstance := kafkago.NewSegmentio()

	// creating kafka topics
	conn, controllerConn, err := segmentioInstance.CreateTopics(cfg.KafkaConfig.BrokerAddress)
	if err != nil {
		log.Fatalf("Unable to create kafka topics | broker address: %s | error: %v", cfg.KafkaConfig.BrokerAddress, err)
	} else {
		log.Println("Successfully created kafka topics!")
	}
	defer conn.Close()
	defer controllerConn.Close()

	db := postgres_conn.ConnectToDB(*cfg)
	defer db.Close()

	// initialise grpc client connections
	grpcClient := grpcclient.NewGRPCClient(*cfg)
	defer grpcClient.InventoryClient().Close()

	// initialise grpc order server
	orderRepo := outbound.NewRepository(db, grpcClient.InventoryClient(), segmentioInstance)
	orderService := services.NewService(*cfg, orderRepo)

	app := grpcserver.Application{
		OrderService: orderService,
		Config:       *cfg,
	}

	go app.InitiateGRPCServer()

	// register authentication microservice with service discovery consul
	serviceDiscovery := consul.NewConsul(*cfg)
	if err := serviceDiscovery.RegisterService(); err != nil {
		log.Fatalf("failed to register authentication microservice with error: %v\n", err)
	}

	// initiate event bus with inventory microservice
	events := services.NewServiceEvents(orderRepo)
	events.ConsumeUpdateInventoryEvent(cfg.KafkaConfig.BrokerAddress, kafkago.TOPIC_UPDATE_ORDER_STATUS)

	select {}
}
