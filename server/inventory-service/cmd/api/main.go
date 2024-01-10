package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	inventory "github.com/LeonLow97/internal"
	"github.com/LeonLow97/pkg/kafkago"
	pb "github.com/LeonLow97/proto"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

var inventoryServicePort = os.Getenv("SERVICE_PORT")

const (
	topicDecrementInventory = "update-inventory-count"
	topicUpdateOrderStatus  = "update-order-status"
	brokerAddress           = "broker:9092"
)

type application struct {
}

func main() {
	app := application{}

	db := app.connectToDB()
	defer db.Close()

	// initiate kafka-go segmentio instance
	segmentioInstance := kafkago.NewSegmentio()

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

	inventoryRepo := inventory.NewRepository(db)
	inventoryService := inventory.NewService(inventoryRepo, segmentioInstance)
	inventory.NewInventoryGRPCHandler(inventoryService)

	go app.initiateGRPCServer(db, segmentioInstance)

	inventoryService.ConsumeKafkaUpdateInventoryCount()

	select {}
}

func (app *application) initiateGRPCServer(db *sql.DB, segmentioInstance *kafkago.Segmentio) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", inventoryServicePort))
	if err != nil {
		log.Fatalf("Failed to start the grpc server with error: %v", err)
	}

	inventoryService := inventory.NewService(inventory.NewRepository(db), segmentioInstance)

	// creates a new grpc server
	grpcServer := grpc.NewServer()
	inventoryServiceServer := inventory.NewInventoryGRPCHandler(inventoryService)

	pb.RegisterInventoryServiceServer(grpcServer, inventoryServiceServer)
	log.Printf("Started inventory gRPC server at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start the inventory gRPC server with error %v", err)
	}
}

func (app *application) connectToDB() *sql.DB {
	// MySQL DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	// open a connection to MySQL database
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening connection to mysql database in inventory service", err)
	}

	// ping mysql database
	if err = conn.Ping(); err != nil {
		log.Fatal("Error pinging mysql database", err)
	}

	return conn
}
