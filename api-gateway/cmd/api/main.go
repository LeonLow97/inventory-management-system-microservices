package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LeonLow97/internal/adapters/outbound"
	"github.com/LeonLow97/internal/core/ports"
	"github.com/LeonLow97/internal/core/usecases"
	"github.com/LeonLow97/internal/delivery/http/handlers"
	"github.com/LeonLow97/internal/delivery/http/ratelimit"
	"github.com/LeonLow97/internal/infrastructure/cache"
	"github.com/LeonLow97/internal/infrastructure/config"
	grpcauthify "github.com/LeonLow97/internal/infrastructure/grpc/authify"
	"github.com/LeonLow97/internal/observability/logging"
	"github.com/LeonLow97/internal/observability/metrics"
)

func main() {
	logger := logging.New(os.Stdout)
	if err := run(logger); err != nil {
		logger.WithError(err).Error("API Gateway stopped with an error")
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

	authifyGrpcClient, err := grpcauthify.NewClient(*cfg)
	if err != nil {
		return fmt.Errorf("create Authify gRPC client: %w", err)
	}
	defer func() {
		if err := authifyGrpcClient.Close(); err != nil {
			logger.WithError(err).Warn("failed to close Authify gRPC client")
		}
	}()

	metrics.InitMetrics()

	rateLimiter, err := ratelimit.NewRateLimiter(ratelimit.AuthRateLimitingConfig{
		Login: ratelimit.LoginRateLimitingConfig{
			Burst: ratelimit.LoginRateLimitConfig{
				Key:         cfg.AuthRateLimit.Login.Burst.Key,
				Window:      cfg.AuthRateLimit.Login.Burst.RateLimit.Window,
				MaxRequests: cfg.AuthRateLimit.Login.Burst.RateLimit.MaxRequests,
			},
			Fail: ratelimit.LoginRateLimitConfig{
				Key:         cfg.AuthRateLimit.Login.Fail.Key,
				Window:      cfg.AuthRateLimit.Login.Fail.RateLimit.Window,
				MaxRequests: cfg.AuthRateLimit.Login.Fail.RateLimit.MaxRequests,
			},
		},
	}, appCache, logger)
	if err != nil {
		return fmt.Errorf("configure login rate limiter: %w", err)
	}

	healthHandler := handlers.NewHealthHandler(dependencyReadinessChecker{
		redis:         appCache,
		authifyHealth: authifyGrpcClient.HealthService(),
	})

	userOutbound := outbound.NewUserOutbound(authifyGrpcClient.AuthifyService())
	userService := usecases.NewUserService(userOutbound, rateLimiter, logger)
	userHandler := handlers.NewUserHandler(handlers.JWTConfig{
		Name:     cfg.AuthJWTToken.Name,
		MaxAge:   cfg.AuthJWTToken.MaxAge,
		Path:     cfg.AuthJWTToken.Path,
		Domain:   cfg.AuthJWTToken.Domain,
		Secure:   cfg.AuthJWTToken.Secure,
		HTTPOnly: cfg.AuthJWTToken.HTTPOnly,
	}, userService, logger)

	app := application{
		cfg:           cfg,
		HealthHandler: healthHandler,
		UserHandler:   userHandler,
	}
	router := app.registerRoutes(logger)

	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.BaseUrl, cfg.Server.Port)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	serverErrors := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- fmt.Errorf("HTTP server failed: %w", err)
		}
	}()

	logger.Info("API Gateway started", ports.LogField{Key: "address", Value: serverAddr})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case receivedSignal := <-quit:
		logger.Info("shutdown signal received", ports.LogField{Key: "signal", Value: receivedSignal.String()})
	case err := <-serverErrors:
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("graceful HTTP shutdown failed: %w", err)
	}

	logger.Info("API Gateway stopped")
	return nil
}
