package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	coreerrors "github.com/LeonLow97/internal/core/errors"
	coreports "github.com/LeonLow97/internal/core/ports"
)

func (r *rateLimiter) LoginRateLimiter(ctx context.Context, email, ip string) error {
	// Burst Rate Limiter
	if err := r.loginBurstRateLimiter(ctx, email, ip); err != nil {
		if isRateLimitExceeded(err) {
			return err
		}
		return fmt.Errorf("login burst rate-limit check failed: %w", err)
	}

	// Fail Rate Limiter
	if err := r.checkLoginFailRateLimiter(ctx, email, ip); err != nil {
		if isRateLimitExceeded(err) {
			return err
		}
		return fmt.Errorf("login failure rate-limit check failed: %w", err)
	}

	return nil
}

func (r *rateLimiter) ResetLoginFailRateLimiter(ctx context.Context, email, ip string) error {
	cfg := r.cfg.Login.Fail
	key := fmt.Sprintf(cfg.Key, email, ip)
	if err := r.store.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to reset login fail rate limiter: %w", err)
	}
	return nil
}

func (r *rateLimiter) IncrementLoginFailRateLimiter(ctx context.Context, email, ip string) error {
	cfg := r.cfg.Login.Fail
	key := fmt.Sprintf(cfg.Key, email, ip)
	count, err := r.fixedWindowCounter(ctx, key, cfg.Window)
	if err != nil {
		return fmt.Errorf("failed to increment login fail rate limiter: %w", err)
	}

	r.logger.WithContext(ctx).Info("login failure counter incremented",
		coreports.LogField{Key: "email", Value: email},
		coreports.LogField{Key: "client_ip", Value: ip},
		coreports.LogField{Key: "count", Value: count},
	)
	return nil
}

func (r *rateLimiter) loginBurstRateLimiter(ctx context.Context, email, ip string) error {
	cfg := r.cfg.Login.Burst
	key := fmt.Sprintf(cfg.Key, email, ip)
	count, err := r.fixedWindowCounter(ctx, key, cfg.Window)
	if err != nil {
		return fmt.Errorf("failed to apply login burst rate limiter: %w", err)
	}

	if count > cfg.MaxRequests {
		return coreerrors.New(coreerrors.KindResourceExhausted, "Rate limit exceeded", nil)
	}
	return nil
}

func (r *rateLimiter) checkLoginFailRateLimiter(ctx context.Context, email, ip string) error {
	cfg := r.cfg.Login.Fail
	key := fmt.Sprintf(cfg.Key, email, ip)

	count, found, err := r.store.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to get login fail counter: %w", err)
	}
	if !found {
		return nil
	}

	countInt, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid login fail counter %q: %w", count, err)
	}
	if countInt >= cfg.MaxRequests {
		return coreerrors.New(coreerrors.KindResourceExhausted, "Rate limit exceeded", nil)
	}
	return nil
}

func isRateLimitExceeded(err error) bool {
	var appErr *coreerrors.Error
	return errors.As(err, &appErr) && appErr.Kind == coreerrors.KindResourceExhausted
}
