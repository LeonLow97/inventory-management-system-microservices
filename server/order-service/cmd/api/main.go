package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"time"

	order "github.com/LeonLow97/internal"
	"github.com/LeonLow97/pkg/aws"
	"github.com/LeonLow97/pkg/config"
	kafkago "github.com/LeonLow97/pkg/kafkago"
	s3client "github.com/LeonLow97/pkg/s3"
	pb "github.com/LeonLow97/proto"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var orderServicePort = os.Getenv("SERVICE_PORT")

type application struct {
	orderService *order.OrderGRPCServer
}

type grpcClientConn struct {
	inventoryConn *grpc.ClientConn
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load config", err)
	}

	app := application{}

	db, err := app.connectToDB()
	if err != nil {
		log.Fatalf("failed to connect to postgres db: %v\n", err)
	}

	// initialize session with aws
	awsSession, err := aws.NewSession(cfg)
	if err != nil {
		log.Fatalln("error getting aws session", err)
	}

	// initialize session with s3
	s3Session := s3client.NewS3(awsSession, 10*time.Second)

	// fileContent := `This is a test file generated from Golang by Jie Wei!`
	// reader := strings.NewReader(fileContent)

	// // test s3 upload object
	// fmt.Println("Bucket name", cfg.AWS.Bucket)
	// loc, err := s3Session.UploadObject(context.Background(), cfg.AWS.Bucket, "/test/temp.txt", reader)
	// if err != nil {
	// 	log.Fatalln("error uploading object to s3", err)
	// }
	// fmt.Println("location", loc)

	// initiate kafka-go segmentio instance
	segmentioInstance := kafkago.NewSegmentio()

	// add update inventory count topic to kafka
	kafkaConfigUpdateInventoryCount := kafkago.NewKafkaConfig("broker:9092", "update-inventory-count")
	conn, controllerConn, err := segmentioInstance.CreateTopics(kafkaConfigUpdateInventoryCount.BrokerAddress)
	if err != nil {
		log.Fatalln("Unable to create kafka topics", err)
	} else {
		log.Println("Successfully created kafka topics!")
	}
	defer conn.Close()
	defer controllerConn.Close()

	grpcClients := app.initiateGRPCClients()
	defer grpcClients.inventoryConn.Close()

	app.setupDBDependencies(db, segmentioInstance, grpcClients, kafkaConfigUpdateInventoryCount, s3Session)

	// running grpc server in the background
	go app.initiateGRPCServer(db, segmentioInstance, kafkaConfigUpdateInventoryCount)

	select {}
}

func (app *application) initiateGRPCClients() *grpcClientConn {
	inventoryConn, err := grpc.Dial(INVENTORY_SERVICE_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("error dialing inventory microservice grpc", err)
	}

	return &grpcClientConn{
		inventoryConn: inventoryConn,
	}
}

func (app *application) initiateGRPCServer(db *sqlx.DB, segmentio *kafkago.Segmentio, kafkaconfig *kafkago.KafkaConfig) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", orderServicePort))
	if err != nil {
		log.Fatalf("Failed to start grpc server with error: %v\n", err)
	}

	// create new grpc server
	grpcServer := grpc.NewServer()

	pb.RegisterOrderServiceServer(grpcServer, app.orderService)
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
