package grpcclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/LeonLow97/internal/config"
	"github.com/LeonLow97/internal/pkg/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient interface {
	OrderClient() *grpc.ClientConn
	InventoryClient() *grpc.ClientConn
	AuthenticationClient() *grpc.ClientConn
}

type grpcClientConn struct {
	orderConn     *grpc.ClientConn
	inventoryConn *grpc.ClientConn
	authConn      *grpc.ClientConn
}

func NewGRPCClient(cfg config.Config, consul *consul.Consul) GRPCClient {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// services := consul.GetServices()
	// authService := services[cfg.AuthService.Name]
	// inventoryService := services[cfg.InventoryService.Name]
	// orderService := services[cfg.OrderService.Name]

	// authConn, err := createGRPCConnection(ctx, fmt.Sprintf("%s:%d", authService.Service, authService.Port))
	authConn, err := createGRPCConnection(ctx, "authentication-service:50051")
	if err != nil {
		log.Fatalf("Error dialing authentication microservice gRPC: %v", err)
	}
	// inventoryConn, err := createGRPCConnection(ctx, fmt.Sprintf("%s:%d", inventoryService.Service, inventoryService.Port))
	// if err != nil {
	// 	log.Fatalf("Error dialing inventory microservice gRPC: %v", err)
	// }
	// orderConn, err := createGRPCConnection(ctx, fmt.Sprintf("%s:%d", orderService.Service, orderService.Port))
	// if err != nil {
	// 	log.Fatalf("Error dialing order microservice gRPC: %v", err)
	// }

	return &grpcClientConn{
		// orderConn:     orderConn,
		// inventoryConn: inventoryConn,
		authConn: authConn,
	}
}

func (c *grpcClientConn) OrderClient() *grpc.ClientConn {
	return c.orderConn
}

func (c *grpcClientConn) InventoryClient() *grpc.ClientConn {
	return c.inventoryConn
}

func (c *grpcClientConn) AuthenticationClient() *grpc.ClientConn {
	return c.authConn
}

func createGRPCConnection(ctx context.Context, serviceURL string) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, serviceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial service at %s: %v", serviceURL, err)
	}
	return conn, nil
}
