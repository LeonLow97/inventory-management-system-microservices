package ratelimit

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/LeonLow97/internal/core/ports"
)

type Store interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Delete(ctx context.Context, key string) error
	IncrementWithExpiry(ctx context.Context, key string, window time.Duration) (int64, error)
}

// AuthRateLimitingConfig contains the login rate-limit policy used by the gateway.
type AuthRateLimitingConfig struct {
	Login LoginRateLimitingConfig
}

// LoginRateLimitingConfig groups the burst and failed-login limits.
type LoginRateLimitingConfig struct {
	Burst LoginRateLimitConfig
	Fail  LoginRateLimitConfig
}

// LoginRateLimitConfig defines one login rate limit.
type LoginRateLimitConfig struct {
	Key         string
	Window      time.Duration
	MaxRequests int64
}

type rateLimiter struct {
	cfg    AuthRateLimitingConfig
	store  Store
	logger ports.Logger
}

func NewRateLimiter(cfg AuthRateLimitingConfig, store Store, logger ports.Logger) (ports.LoginAttemptLimiter, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return &rateLimiter{
		cfg:    cfg,
		store:  store,
		logger: logger,
	}, nil
}

func validateConfig(cfg AuthRateLimitingConfig) error {
	policies := []struct {
		name   string
		config LoginRateLimitConfig
	}{
		{name: "burst", config: cfg.Login.Burst},
		{name: "failure", config: cfg.Login.Fail},
	}

	for _, policy := range policies {
		if strings.TrimSpace(policy.config.Key) == "" || strings.Count(policy.config.Key, "%s") != 2 {
			return fmt.Errorf("invalid %s login rate limit: key must contain two %%s placeholders", policy.name)
		}
		if policy.config.Window < time.Millisecond {
			return fmt.Errorf("invalid %s login rate limit: window must be at least 1ms", policy.name)
		}
		if policy.config.MaxRequests <= 0 {
			return fmt.Errorf("invalid %s login rate limit: max requests must be positive", policy.name)
		}
	}
	return nil
}
