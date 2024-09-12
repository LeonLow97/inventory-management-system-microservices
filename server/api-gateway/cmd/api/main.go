package main

import (
	"fmt"
	"os"

	"github.com/LeonLow97/internal/adapters/inbound/web"
	grpcclient "github.com/LeonLow97/internal/adapters/outbound/grpc"
	"github.com/LeonLow97/internal/config"
	"github.com/LeonLow97/internal/core/services/auth"
	"github.com/LeonLow97/internal/core/services/inventory"
	"github.com/LeonLow97/internal/core/services/order"
	"github.com/LeonLow97/internal/core/services/user"
	"go.uber.org/zap"
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
	// Setup logger using Uber Zap
	setupLog()
	defer logger.Sync() // flushes buffer, if any

	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("failed to load config with error", zap.Error(err))
		return
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
	logger.Info("Starting API Gateway for Inventory Management System!")
	apiGatewayPort := fmt.Sprintf(":%d", cfg.Server.Port)
	if err := router.Run(apiGatewayPort); err != nil {
		logger.Fatal("failed to run server", zap.Error(err))
	}
}

const logPath = "../../logs/gateway.log"

var logger *zap.Logger

func setupLog() {
	// Create or open the log file
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("failed to open log file: %v\n", err)
		return
	}
	defer file.Close()

	// Configure Zap logger
	c := zap.NewProductionConfig()
	c.OutputPaths = []string{"stdout", logPath}
	logger, err = c.Build()
	if err != nil {
		fmt.Printf("failed to build logger: %v\n", err)
		return
	}
}
