package main

import (
	"log"

	"github.com/LeonLow97/internal/adapters/outbound"
	"github.com/LeonLow97/internal/core/services"
	"github.com/LeonLow97/internal/pkg/config"
	"github.com/LeonLow97/internal/pkg/grpcserver"
	"github.com/LeonLow97/internal/pkg/kafkago"
	mysql_conn "github.com/LeonLow97/internal/pkg/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed to load config with error:", err)
	}

	// initiate kafka-go segmentio instance
	segmentioInstance := kafkago.NewSegmentio()

	conn, controllerConn, err := segmentioInstance.CreateTopics("broker:9092")
	if err != nil {
		log.Fatalf("Unable to create kafka topics | broker address: %s | error: %v", cfg.KafkaConfig.BrokerAddress, err)
	} else {
		log.Println("Successfully created kafka topics!")
	}

	defer conn.Close()
	defer controllerConn.Close()

	mySQLConn := mysql_conn.ConnectToMySQL(*cfg)
	defer mySQLConn.Close()

	// initialise grpc inventory server
	inventoryRepo := outbound.NewRepository(mySQLConn, segmentioInstance)
	inventoryService := services.NewService(inventoryRepo)

	app := &grpcserver.Application{
		Service: inventoryService,
		Config:  *cfg,
	}

	go app.InitiateGRPCServer()

	// initiate event bus with order microservice
	events := services.NewServiceEvents(inventoryRepo)
	events.ConsumeUpdateInventoryEvent(cfg.KafkaConfig.BrokerAddress, kafkago.TOPIC_DECREMENT_INVENTORY, kafkago.TOPIC_UPDATE_ORDER_STATUS)

	select {}
}
