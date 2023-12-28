package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/LeonLow97/internal/authenticate"
	"github.com/LeonLow97/internal/users"
	pb "github.com/LeonLow97/proto"
)

var authenticationServicePort = os.Getenv("SERVICE_PORT")

type application struct {

}

func main() {
	app := application{

	}

	db, err := app.connectToDB()
	if err != nil {
		log.Fatalln(err)
	}

	go app.initiateGRPCServer(db)

	r := app.routes(db)
	log.Println("Auth Service is running on port", authenticationServicePort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", "8002"), r); err != nil {
		log.Fatalf("Failed to start authentication microservice with error %v", err)
	}
}

func (app *application) initiateGRPCServer(db *sqlx.DB) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", authenticationServicePort))
	if err != nil {
		log.Fatalf("Failed to start the grpc server with error: %v", err)
	}

	authService := authenticate.NewService(authenticate.NewRepo(db))
	userService := users.NewService(users.NewRepo(db))

	// creates a new grpc server
	grpcServer := grpc.NewServer()
	authServiceServer := authenticate.NewAuthenticateGRPCHandler(authService)
	usersServiceServer := users.NewUsersGRPCHandler(userService)

	pb.RegisterAuthenticationServiceServer(grpcServer, authServiceServer)
	pb.RegisterUserServiceServer(grpcServer, usersServiceServer)
	log.Printf("Started authentication gRPC server at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start authentication gRPC server with error %v", err)
	}
}

func (app *application) connectToDB() (*sqlx.DB, error) {
	// environment := "development"
	// var databaseURL string

	// config, err := config.LoadConfig(environment)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }

	// if environment == "development" {
	// 	databaseURL = config.Development.URL
	// }

	// databaseURL := "postgres://postgres:password@localhost:5432/imsdb?sslmode=disable"

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

func runMigrations(migrationsPath string, db *sql.DB) error {
	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"imsdb",
		instance,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
