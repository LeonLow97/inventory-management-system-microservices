package main

import (
	"log"
	"os"

	"github.com/LeonLow97/internal/adapters/outbound"
	"github.com/LeonLow97/internal/core/services"
	"github.com/LeonLow97/internal/pkg/kafkago"
	mysql_conn "github.com/LeonLow97/internal/pkg/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var inventoryServicePort = os.Getenv("SERVICE_PORT")

const (
	topicDecrementInventory = "update-inventory-count" // consume
	topicUpdateOrderStatus  = "update-order-status"    // produce
	brokerAddress           = "broker:9092"
)

type application struct {
	service services.Service
}

func main() {
	// initiate kafka-go segmentio instance
	segmentioInstance := kafkago.NewSegmentio()

	segmentioInstance.AddTopicConfig(topicDecrementInventory, 1, 1)
	segmentioInstance.AddTopicConfig(topicUpdateOrderStatus, 1, 1)
	conn, controllerConn, err := segmentioInstance.CreateTopics(brokerAddress)
	if err != nil {
		log.Fatalln("Unable to create kafka topics", err)
	}
	log.Println("Successfully created kafka topics!")
	defer conn.Close()
	defer controllerConn.Close()

	db := mysql_conn.ConnectToMySQL()
	defer db.Close()

	// initialise grpc inventory server
	inventoryRepo := outbound.NewRepository(db, segmentioInstance)
	inventoryService := services.NewService(inventoryRepo)

	app := &application{
		service: inventoryService,
	}

	go app.InitiateGRPCServer()

	// initiate event bus with order microservice
	events := services.NewServiceEvents(inventoryRepo)
	events.ConsumeUpdateInventoryEvent(brokerAddress, topicDecrementInventory, topicUpdateOrderStatus)

	select {}
}
