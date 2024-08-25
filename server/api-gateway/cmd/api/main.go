package main

import (
	"fmt"
	"log"

	"github.com/LeonLow97/internal/adapters/inbound/web"
	grpcclient "github.com/LeonLow97/internal/adapters/outbound/grpc"
	"github.com/LeonLow97/internal/config"
	"github.com/LeonLow97/internal/core/services/auth"
	"github.com/LeonLow97/internal/core/services/inventory"
	"github.com/LeonLow97/internal/core/services/order"
	"github.com/LeonLow97/internal/core/services/user"
)

type application struct {
	Config           *config.Config
	GRPCClient       grpcclient.GRPCClient
	AuthHandler      *web.AuthHandler
	UserHandler      *web.UserHandler
	InventoryHandler *web.InventoryHandler
	OrderHandler     *web.OrderHandler
}

func main() {
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v\n", err)
	}

	grpcClient := grpcclient.NewGRPCClient(cfg)
	defer grpcClient.AuthenticationClient().Close()
	defer grpcClient.InventoryClient().Close()
	defer grpcClient.OrderClient().Close()

	// instantiating auth microservice
	authRepo := grpcclient.NewAuthRepo(grpcClient.AuthenticationClient())
	authService := auth.NewAuthService(authRepo)
	authHandler := web.NewAuthHandler(authService)

	// instantiating user microservice
	userRepo := grpcclient.NewUserRepo(grpcClient.AuthenticationClient())
	userService := user.NewUserService(userRepo)
	userHandler := web.NewUserHandler(userService)

	// instantiating inventory microservice
	inventoryRepo := grpcclient.NewInventoryRepo(grpcClient.InventoryClient())
	inventoryService := inventory.NewInventoryService(inventoryRepo)
	inventoryHandler := web.NewInventoryHandler(inventoryService)

	// instantiating order microservice
	orderRepo := grpcclient.NewOrderRepo(grpcClient.OrderClient())
	orderService := order.NewOrderService(orderRepo)
	orderHandler := web.NewOrderHandler(orderService)

	// setup application config
	app := application{
		Config:           cfg,
		GRPCClient:       grpcClient,
		AuthHandler:      authHandler,
		UserHandler:      userHandler,
		InventoryHandler: inventoryHandler,
		OrderHandler:     orderHandler,
	}

	// getting router with gin engine
	router := app.routes()

	// Using gin to start api gateway server, exit status 1 if fail to start server
	apiGatewayPort := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting API Gateway on Port %d", cfg.Server.Port)
	if err := router.Run(apiGatewayPort); err != nil {
		log.Fatal(err)
	}
}
