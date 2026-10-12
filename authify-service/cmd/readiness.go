package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type databaseReadinessChecker interface {
	PingContext(context.Context) error
}

type redisReadinessChecker interface {
	Ping(context.Context) error
}

type dependencyReadinessChecker struct {
	database databaseReadinessChecker
	redis    redisReadinessChecker
}

func (c dependencyReadinessChecker) Check(ctx context.Context) error {
	if err := c.database.PingContext(ctx); err != nil {
		return fmt.Errorf("PostgreSQL is unavailable: %w", err)
	}
	if err := c.redis.Ping(ctx); err != nil {
		return fmt.Errorf("redis is unavailable: %w", err)
	}
	return nil
}

type dependencyHealthServer struct {
	*health.Server
	readiness dependencyReadinessChecker
}

func (s dependencyHealthServer) Check(ctx context.Context, request *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	if err := s.readiness.Check(ctx); err != nil {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
	return s.Server.Check(ctx, request)
}
