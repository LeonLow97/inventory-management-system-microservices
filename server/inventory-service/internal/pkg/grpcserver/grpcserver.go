package grpcserver

import (
	"fmt"
	"log"
	"net"

	inbound "github.com/LeonLow97/internal/adapters/inbound/grpcserver"
	"github.com/LeonLow97/internal/core/services"
	"github.com/LeonLow97/internal/pkg/config"
	pb "github.com/LeonLow97/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Application struct {
	Service services.Service
	Config  config.Config
}

func (app *Application) InitiateGRPCServer() {
	// creates a new grpc server
	grpcServer := grpc.NewServer()

	// initialize health server for grpc
	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthService)
	healthService.SetServingStatus(app.Config.Server.Name, grpc_health_v1.HealthCheckResponse_SERVING)

	// gRPC reflection for service discovery by grpc clients
	// allows gRPC clients to discover the services and methods exposed by a gRPC server at runtime
	// useful for service discovery, introspection tools (lik grpcurl), and debugging
	reflection.Register(grpcServer)

	inventoryGRPCServer := inbound.NewInventoryGRPCServer(app.Service)
	pb.RegisterInventoryServiceServer(grpcServer, inventoryGRPCServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", app.Config.Server.Port))
	if err != nil {
		log.Fatalf("Failed to start the grpc server with error: %v", err)
	}
	log.Printf("Started inventory gRPC server at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start the inventory gRPC server with error %v", err)
	}
}
