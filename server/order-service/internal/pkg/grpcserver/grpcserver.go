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

// Ensuring gRPC Server Health Checks on Kubernetes: A Comprehensive Guide
// Link: https://medium.com/@blackhorseya/ensuring-grpc-server-health-checks-on-kubernetes-a-comprehensive-guide-86aac08ad5b0

type Application struct {
	OrderService services.Service
	Config       config.Config
}

func (app *Application) InitiateGRPCServer() {
	// creates a new grpc server
	grpcServer := grpc.NewServer()

	// initialize health server for grpc
	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthService)
	healthService.SetServingStatus(app.Config.Server.Name, grpc_health_v1.HealthCheckResponse_SERVING) // to set the status of grpc server when performing health checks

	// gRPC reflection for service discovery by grpc clients
	// allows gRPC clients to discover the services and methods exposed by a gRPC server at runtime
	// useful for service discovery, introspection tools (lik grpcurl), and debugging
	reflection.Register(grpcServer)

	// register order grpc server
	orderGRPCServer := inbound.NewOrderGRPCServer(app.OrderService)
	pb.RegisterOrderServiceServer(grpcServer, orderGRPCServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", app.Config.Server.Port))
	if err != nil {
		log.Fatalf("Failed to start grpc server with error: %v\n", err)
	}
	log.Printf("Started order gRPC server at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start order gRPC server with error %v", err)
	}
}
