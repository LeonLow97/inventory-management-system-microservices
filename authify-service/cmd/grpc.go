package main

import (
	pb "github.com/LeonLow97/authify/grpc-proto/pb"
	"github.com/LeonLow97/internal/core/ports"
	"github.com/LeonLow97/internal/core/usecases"
	"github.com/LeonLow97/internal/delivery/grpc/handlers"
	"github.com/LeonLow97/internal/delivery/grpc/interceptor"
	"github.com/LeonLow97/internal/infrastructure/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type application struct {
	cfg       config.Config
	service   usecases.Service
	logger    ports.Logger
	readiness dependencyReadinessChecker
}

func (app *application) newGRPCServer() *grpc.Server {
	server := grpc.NewServer(
		grpc.MaxConcurrentStreams(1000),
		grpc.ChainUnaryInterceptor(
			interceptor.RequestID(),
			interceptor.RequestLogger(app.logger),
			interceptor.Recovery(app.logger),
		),
	)

	app.registerHealth(server)
	app.registerReflection(server)
	app.registerAuthifyService(server)
	return server
}

func (app *application) registerHealth(server *grpc.Server) {
	healthService := health.NewServer()
	healthService.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthService.SetServingStatus(app.cfg.Server.Name, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, dependencyHealthServer{
		Server:    healthService,
		readiness: app.readiness,
	})
}

func (app *application) registerReflection(server *grpc.Server) {
	reflection.Register(server)
}

func (app *application) registerAuthifyService(server *grpc.Server) {
	authifyServer := handlers.NewUserHandler(app.service)
	pb.RegisterAuthifyServiceServer(server, authifyServer)
}
