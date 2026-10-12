package authify

import (
	"fmt"

	pb "github.com/LeonLow97/authify/grpc-proto/pb"
	"github.com/LeonLow97/internal/infrastructure/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Client interface {
	AuthifyService() pb.AuthifyServiceClient
	HealthService() grpc_health_v1.HealthClient
	Close() error
}

type client struct {
	conn          *grpc.ClientConn
	authifyClient pb.AuthifyServiceClient
}

// NewClient connects to the authify service over gRPC.
func NewClient(cfg config.Config) (Client, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(10 << 20), // 10MB
		),
	}

	serverAddr := fmt.Sprintf("%s:%d", cfg.AuthifyService.BaseUrl, cfg.AuthifyService.Port)
	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial auth service %q: %w", serverAddr, err)
	}

	return &client{
		conn:          conn,
		authifyClient: pb.NewAuthifyServiceClient(conn),
	}, nil
}

func (c *client) AuthifyService() pb.AuthifyServiceClient {
	return c.authifyClient
}

func (c *client) HealthService() grpc_health_v1.HealthClient {
	return grpc_health_v1.NewHealthClient(c.conn)
}

func (c *client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
