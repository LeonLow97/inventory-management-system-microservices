package ratelimit

import (
	"context"
	"fmt"
	"time"
)

// fixedWindowCounter increments a counter and sets its expiry atomically.
func (r *rateLimiter) fixedWindowCounter(
	ctx context.Context,
	key string,
	window time.Duration,
) (int64, error) {
	count, err := r.store.IncrementWithExpiry(ctx, key, window)
	if err != nil {
		return 0, fmt.Errorf("failed to increment fixed-window counter %q: %w", key, err)
	}
	return count, nil
}
