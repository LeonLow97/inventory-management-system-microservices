package grpcclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/LeonLow97/internal/config"
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

func NewGRPCClient(cfg *config.Config) GRPCClient {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	authConn, err := createGRPCConnection(ctx, cfg.AuthService.URL)
	if err != nil {
		log.Fatalf("Error dialing authentication microservice gRPC: %v", err)
	}
	inventoryConn, err := createGRPCConnection(ctx, cfg.InventoryService.URL)
	if err != nil {
		log.Fatalf("Error dialing inventory microservice gRPC: %v", err)
	}
	orderConn, err := createGRPCConnection(ctx, cfg.OrderService.URL)
	if err != nil {
		log.Fatalf("Error dialing order microservice gRPC: %v", err)
	}

	return &grpcClientConn{
		orderConn:     orderConn,
		inventoryConn: inventoryConn,
		authConn:      authConn,
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
