package main

import (
	"context"
	"fmt"
	"time"

	"log"

	"github.com/LeonLow97/internal/adapters/inbound/web"
	grpcclient "github.com/LeonLow97/internal/adapters/outbound/grpc"
	"github.com/LeonLow97/internal/config"
	"github.com/LeonLow97/internal/core/services/auth"
	"github.com/LeonLow97/internal/core/services/user"
	"github.com/LeonLow97/internal/pkg/cache"
	"github.com/LeonLow97/internal/pkg/consul"
)

type application struct {
	Config           *config.Config
	GRPCClient       grpcclient.GRPCClient
	AppCache         cache.Cache
	AuthHandler      *web.AuthHandler
	UserHandler      *web.UserHandler
	InventoryHandler *web.InventoryHandler
	OrderHandler     *web.OrderHandler
}

func main() {
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed to load config with error", err)
		return
	}

	// Load application cache
	appCache := cache.NewRedisClient(*cfg)

	// create a consul client
	hashicorpConsul := consul.NewConsul(*cfg)
	hashicorpClient, err := hashicorpConsul.NewConsul(*cfg)
	if err != nil {
		log.Fatalf("failed to create hashicorp consul client with error: %v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := hashicorpClient.RefreshServices(ctx); err != nil {
		log.Fatalln("failed to refresh services")
	}

	grpcClient := grpcclient.NewGRPCClient(*cfg, hashicorpClient)
	defer grpcClient.AuthenticationClient().Close()
	// defer grpcClient.InventoryClient().Close()
	// defer grpcClient.OrderClient().Close()

	// instantiating auth microservice
	authRepo := grpcclient.NewAuthRepo(grpcClient.AuthenticationClient())
	authService := auth.NewAuthService(authRepo)
	authHandler := web.NewAuthHandler(authService)

	// instantiating user microservice
	userRepo := grpcclient.NewUserRepo(grpcClient.AuthenticationClient())
	userService := user.NewUserService(userRepo)
	userHandler := web.NewUserHandler(userService)

	// // instantiating inventory microservice
	// inventoryRepo := grpcclient.NewInventoryRepo(grpcClient.InventoryClient())
	// inventoryService := inventory.NewInventoryService(inventoryRepo)
	// inventoryHandler := web.NewInventoryHandler(inventoryService)

	// // instantiating order microservice
	// orderRepo := grpcclient.NewOrderRepo(grpcClient.OrderClient())
	// orderService := order.NewOrderService(orderRepo)
	// orderHandler := web.NewOrderHandler(orderService)

	// setup application config
	app := application{
		Config:      cfg,
		AppCache:    appCache,
		GRPCClient:  grpcClient,
		AuthHandler: authHandler,
		UserHandler: userHandler,
		// InventoryHandler: inventoryHandler,
		// OrderHandler:     orderHandler,
	}

	// getting router with gin engine
	router := app.routes()

	// Using gin to start api gateway server, exit status 1 if fail to start server
	log.Println("Starting API Gateway for Inventory Management System!")
	apiGatewayPort := fmt.Sprintf(":%d", cfg.Server.Port)
	if err := router.Run(apiGatewayPort); err != nil {
		log.Fatal("failed to run server", err)
	}
}
