package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	inventory "github.com/LeonLow97/internal"
	kafkago "github.com/LeonLow97/internal/kafka"
	pb "github.com/LeonLow97/proto"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

var inventoryServicePort = os.Getenv("SERVICE_PORT")

const (
	topicDecrementInventory = "DECREMENT_INVENTORY"
	brokerAddress           = "broker:9092"
)

type application struct {
}

func main() {
	app := application{}

	db := app.connectToDB()
	defer db.Close()

	go app.initiateGRPCServer(db)

	// initiate kafka-go segmentio instance
	segmentioInstance := kafkago.NewSegmentio()

	segmentioInstance.AddTopicConfig(topicDecrementInventory, 1, 1)
	conn, controllerConn, err := segmentioInstance.CreateTopics(brokerAddress)
	if err != nil {
		log.Fatalln("Unable to create kafka topics", err)
	}
	log.Println("Successfully created kafka topics!")
	defer conn.Close()
	defer controllerConn.Close()

	// Consume messages from order management microservice
	go func() {
		if err := segmentioInstance.Consumer(brokerAddress, topicDecrementInventory); err != nil {
			log.Printf("failed to consume message for %s topic: %v\n", topicDecrementInventory, err)
		}
	}()

	app.routes(db)

	select {}
}

func (app *application) initiateGRPCServer(db *sql.DB) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", inventoryServicePort))
	if err != nil {
		log.Fatalf("Failed to start the grpc server with error: %v", err)
	}

	inventoryService := inventory.NewService(inventory.NewRepository(db))

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
