package main

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/LeonLow97/internal/adapters/outbound"
	"github.com/LeonLow97/internal/core/services"
	"github.com/LeonLow97/internal/pkg/config"
	"github.com/LeonLow97/internal/pkg/consul"
	"github.com/LeonLow97/internal/pkg/grpcserver"
	postgres_conn "github.com/LeonLow97/internal/pkg/postgres"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed to load config with error:", err)
	}

	db, err := postgres_conn.ConnectToPostgreSQL(*cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres db: %v\n", err)
	}

	// initialize grpc authentication and user services
	repo := outbound.NewRepository(db)
	service := services.NewService(repo, *cfg)

	app := grpcserver.Application{
		Config:  *cfg,
		Service: service,
	}

	go app.InitiateGRPCServer()

	// register authentication microservice with service discovery consul
	serviceDiscovery := consul.NewConsul(*cfg)
	if err := serviceDiscovery.RegisterService(); err != nil {
		log.Fatalf("failed to register authentication microservice with error: %v\n", err)
	}

	select {}
}
