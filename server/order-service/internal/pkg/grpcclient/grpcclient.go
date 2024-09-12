package grpcclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/LeonLow97/internal/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient interface {
	InventoryClient() *grpc.ClientConn
}

type grpcClientConn struct {
	inventoryConn *grpc.ClientConn
}

func NewGRPCClient(cfg config.Config) GRPCClient {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	inventoryConn, err := createGRPCConnection(ctx, cfg.InventoryService.URL)
	if err != nil {
		log.Fatalf("failed to dial inventory microservice with error: %v\n", err)
	}

	return &grpcClientConn{
		inventoryConn: inventoryConn,
	}
}

func (c *grpcClientConn) InventoryClient() *grpc.ClientConn {
	return c.inventoryConn
}

func createGRPCConnection(ctx context.Context, serviceURL string) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, serviceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial service at %s: %v", serviceURL, err)
	}
	return conn, nil
}
