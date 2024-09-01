package main

import (
	"log"
	"os"

	outbound_kafka "github.com/LeonLow97/internal/adapters/outbound/kafka"
	outbound_mysql "github.com/LeonLow97/internal/adapters/outbound/mysql"
	"github.com/LeonLow97/internal/core/services"
	"github.com/LeonLow97/pkg/kafkago"
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

	db := outbound_mysql.ConnectToMySQL()
	defer db.Close()

	// initialise grpc inventory server
	mySQLAdapter := outbound_mysql.NewMySQLAdapter(db)
	inventoryService := services.NewService(mySQLAdapter)

	app := &application{
		service: inventoryService,
	}

	go app.InitiateGRPCServer(db, segmentioInstance)

	// initiate event bus with order microservice
	eventBusAdapter := outbound_kafka.NewKafkaAdapter(segmentioInstance)
	events := services.NewServiceEvents(mySQLAdapter, eventBusAdapter)
	events.ConsumeUpdateInventoryEvent(brokerAddress, topicDecrementInventory, topicUpdateOrderStatus)

	select {}
}
