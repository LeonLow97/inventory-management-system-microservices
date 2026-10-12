package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc/health/grpc_health_v1"
)

type redisReadinessChecker interface {
	Ping(context.Context) error
}

type dependencyReadinessChecker struct {
	redis         redisReadinessChecker
	authifyHealth grpc_health_v1.HealthClient
}

func (c dependencyReadinessChecker) Check(ctx context.Context) error {
	if err := c.redis.Ping(ctx); err != nil {
		return fmt.Errorf("redis down: %w", err)
	}

	response, err := c.authifyHealth.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return fmt.Errorf("authify service down: %w", err)
	}
	if response.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
		return fmt.Errorf("authify service is not serving: %s", response.GetStatus())
	}

	return nil
}
