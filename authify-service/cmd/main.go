package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/LeonLow97/internal/adapters/outbound/postgres"
	"github.com/LeonLow97/internal/core/ports"
	"github.com/LeonLow97/internal/core/usecases"
	"github.com/LeonLow97/internal/infrastructure/cache"
	"github.com/LeonLow97/internal/infrastructure/config"
	"github.com/LeonLow97/internal/infrastructure/database"
	"github.com/LeonLow97/internal/observability/logging"
	"google.golang.org/grpc"
)

func main() {
	logger := logging.New(os.Stdout)
	if err := run(logger); err != nil {
		logger.WithError(err).Error("Authify service stopped with an error")
		os.Exit(1)
	}
}

func run(logger ports.Logger) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	appCache, err := cache.NewRedisCache(*cfg)
	if err != nil {
		return fmt.Errorf("connect to Redis: %w", err)
	}
	defer func() {
		if err := appCache.Close(); err != nil {
			logger.WithError(err).Warn("failed to close Redis cache")
		}
	}()
	logger.Info("connected to Redis")

	dbClient, err := database.GetDB(*cfg)
	if err != nil {
		return fmt.Errorf("connect to PostgreSQL: %w", err)
	}
	defer func() {
		if err := dbClient.Close(); err != nil {
			logger.WithError(err).Warn("failed to close PostgreSQL connection")
		}
	}()
	logger.Info("connected to PostgreSQL")

	userRepository := postgres.NewUserRepository(dbClient.DB)
	service := usecases.NewService(usecases.Config{
		JWTConfig: usecases.JWTConfig{
			SecretKey: cfg.JWTConfig.SecretKey,
			Expiry:    cfg.JWTConfig.Expiry,
		},
	}, userRepository, appCache, logger)

	app := application{
		cfg:       *cfg,
		service:   service,
		logger:    logger,
		readiness: dependencyReadinessChecker{database: dbClient, redis: appCache},
	}
	grpcServer := app.newGRPCServer()
	defer grpcServer.Stop()

	address := fmt.Sprintf(":%d", cfg.Server.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen for gRPC on %s: %w", address, err)
	}

	serverErrors := make(chan error, 1)
	go func() {
		if err := grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			serverErrors <- fmt.Errorf("gRPC server failed: %w", err)
		}
	}()

	logger.Info("Authify gRPC server started", ports.LogField{Key: "address", Value: address})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case receivedSignal := <-quit:
		logger.Info("shutdown signal received", ports.LogField{Key: "signal", Value: receivedSignal.String()})
	case err := <-serverErrors:
		return err
	}

	gracefulStop := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(gracefulStop)
	}()

	shutdownTimer := time.NewTimer(5 * time.Second)
	defer shutdownTimer.Stop()
	select {
	case <-gracefulStop:
	case <-shutdownTimer.C:
		logger.Warn("gRPC graceful shutdown timed out")
		grpcServer.Stop()
		<-gracefulStop
	}

	logger.Info("Authify service stopped")
	return nil
}
