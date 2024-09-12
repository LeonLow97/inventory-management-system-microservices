package main

import (
	"log"

	"github.com/LeonLow97/internal/adapters/outbound"
	"github.com/LeonLow97/internal/core/services"
	"github.com/LeonLow97/internal/pkg/config"
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
		log.Fatalln("Unable to create kafka topics", err)
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

	// initiate event bus with inventory microservice
	events := services.NewServiceEvents(orderRepo)
	events.ConsumeUpdateInventoryEvent(cfg.KafkaConfig.BrokerAddress, kafkago.TOPIC_UPDATE_ORDER_STATUS)

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
