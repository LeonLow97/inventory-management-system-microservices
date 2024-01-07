package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"

	order "github.com/LeonLow97/internal"
	kafkago "github.com/LeonLow97/internal/kafkago"
	pb "github.com/LeonLow97/proto"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

var orderServicePort = os.Getenv("SERVICE_PORT")

type application struct {
}

func main() {
	app := application{}

	db, err := app.connectToDB()
	if err != nil {
		log.Fatalf("failed to connect to postgres db: %v\n", err)
	}

	app.setupDBDependencies(db)

	// initiate kafka-go segmentio instance
	segmentioInstance := kafkago.NewSegmentio()

	// add update inventory count topic to kafka
	kafkaConfigUpdateInventoryCount := kafkago.NewKafkaConfig("broker:9092", "update-inventory-count")
	segmentioInstance.AddTopicConfig(kafkaConfigUpdateInventoryCount.TopicName, 1, 1)
	conn, controllerConn, err := segmentioInstance.CreateTopics(kafkaConfigUpdateInventoryCount.BrokerAddress)
	if err != nil {
		log.Fatalln("Unable to create kafka topics", err)
	} else {
		log.Println("Successfully created kafka topics!")
	}
	defer conn.Close()
	defer controllerConn.Close()

	// running grpc server in the background
	go app.initiateGRPCServer(db, segmentioInstance, kafkaConfigUpdateInventoryCount)

	select {}
}

func (app *application) initiateGRPCServer(db *sqlx.DB, segmentio *kafkago.Segmentio, kafkaconfig *kafkago.KafkaConfig) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", orderServicePort))
	if err != nil {
		log.Fatalf("Failed to start grpc server with error: %v\n", err)
	}

	orderRepo := order.NewRepository(db)
	orderService := order.NewService(orderRepo, segmentio, kafkaconfig)

	// create new grpc server
	grpcServer := grpc.NewServer()
	orderServiceServer := order.NewOrderGRPCHandler(orderService)

	pb.RegisterOrderServiceServer(grpcServer, orderServiceServer)
	log.Printf("Started order gRPC server at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start order gRPC server with error %v", err)
	}
}

func (app *application) connectToDB() (*sqlx.DB, error) {
	// Construct the DSN string based on environment variables
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	connConfig, err := pgx.ParseConfig(databaseURL)
	if err != nil {
		errMsg := err.Error()
		errMsg = regexp.MustCompile(`(://[^:]+:).+(@.+)`).ReplaceAllString(errMsg, "$1*****$2")
		errMsg = regexp.MustCompile(`(password=).+(\s+)`).ReplaceAllString(errMsg, "$1*****$2")
		return nil, fmt.Errorf("parsing DSN failed: %s", errMsg)
	}
	connStr := stdlib.RegisterConnConfig(connConfig)
	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// if err := runMigrations("migrations", db.DB); err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }

	return db, nil
}
