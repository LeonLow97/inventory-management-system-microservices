package main

import (
	"log"
	"os"

	"github.com/LeonLow97/internal/adapters/outbound"
	"github.com/LeonLow97/internal/core/services"
	"github.com/LeonLow97/pkg/config"
	grpc_conn "github.com/LeonLow97/pkg/grpc"
	kafkago "github.com/LeonLow97/pkg/kafkago"
	postgres_conn "github.com/LeonLow97/pkg/postgres"
)

var orderServicePort = os.Getenv("SERVICE_PORT")

const (
	topicDecrementInventory = "update-inventory-count"
	topicUpdateOrderStatus  = "update-order-status"
	brokerAddress           = "broker:9092"
)

type application struct {
	orderService services.Service
}

func main() {
	// initiate kafka-go segmentio instance
	segmentioInstance := kafkago.NewSegmentio()

	// add update inventory count topic to kafka
	segmentioInstance.AddTopicConfig(topicDecrementInventory, 1, 1)
	segmentioInstance.AddTopicConfig(topicUpdateOrderStatus, 1, 1)
	conn, controllerConn, err := segmentioInstance.CreateTopics(brokerAddress)
	if err != nil {
		log.Fatalln("Unable to create kafka topics", err)
	} else {
		log.Println("Successfully created kafka topics!")
	}
	defer conn.Close()
	defer controllerConn.Close()

	_, err = config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load config", err)
	}

	db := postgres_conn.ConnectToDB()
	defer db.Close()

	// initialise grpc client connections
	grpcClient := grpc_conn.NewGRPCClient()
	defer grpcClient.InventoryClient().Close()

	// initialise grpc order server
	orderRepo := outbound.NewRepository(db, grpcClient.InventoryClient(), segmentioInstance)
	orderService := services.NewService(orderRepo)

	app := application{
		orderService: orderService,
	}

	go app.InitiateGRPCServer()

	// initiate event bus with inventory microservice
	events := services.NewServiceEvents(orderRepo)
	events.ConsumeUpdateInventoryEvent(brokerAddress, topicUpdateOrderStatus)

	select {}

	// // initialize session with aws
	// awsSession, err := aws.NewSession(cfg)
	// if err != nil {
	// 	log.Fatalln("error getting aws session", err)
	// }

	// // initialize session with s3
	// s3Session := s3client.NewS3(awsSession, 10*time.Second)

	// fileContent := `This is a test file generated from Golang by Jie Wei!`
	// reader := strings.NewReader(fileContent)

	// // test s3 upload object
	// fmt.Println("Bucket name", cfg.AWS.Bucket)
	// loc, err := s3Session.UploadObject(context.Background(), cfg.AWS.Bucket, "/test/temp.txt", reader)
	// if err != nil {
	// 	log.Fatalln("error uploading object to s3", err)
	// }
	// fmt.Println("location", loc)
}
