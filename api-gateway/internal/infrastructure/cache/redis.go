package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/LeonLow97/internal/infrastructure/config"
	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client redisCommander
}

type redisCommander interface {
	redis.Scripter
	Ping(ctx context.Context) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Close() error
}

// incrementWithExpiryScript keeps rate-limit counter updates and TTL handling atomic.
// This prevents counters from being left without an expiry.
var incrementWithExpiryScript = redis.NewScript(`
local ttlMilliseconds = tonumber(ARGV[1])
if not ttlMilliseconds or ttlMilliseconds <= 0 then
  return redis.error_reply('rate-limit window must be a positive number of milliseconds')
end

if redis.call('EXISTS', KEYS[1]) == 0 then
  redis.call('SET', KEYS[1], 1, 'PX', ttlMilliseconds)
  return 1
end

if redis.call('PTTL', KEYS[1]) < 0 then
  local expirySet = redis.call('PEXPIRE', KEYS[1], ttlMilliseconds)
  if expirySet ~= 1 then
    return redis.error_reply('failed to set expiry for rate-limit counter')
  end
end

return redis.call('INCR', KEYS[1])
`)

func NewRedisCache(cfg config.Config) (*RedisCache, error) {
	return newRedisCache(redis.NewClient(redisClientOptions(cfg)))
}

func redisClientOptions(cfg config.Config) *redis.Options {
	return &redis.Options{
		Addr:            fmt.Sprintf("%s:%d", cfg.RedisConfig.Address, cfg.RedisConfig.Port),
		Password:        cfg.RedisConfig.Password,
		DB:              cfg.RedisConfig.DatabaseIndex,
		MaxRetries:      3,
		MinRetryBackoff: 100 * time.Millisecond,
		MaxRetryBackoff: 1 * time.Second,
	}
}

func newRedisCache(client redisCommander) (*RedisCache, error) {
	cache := &RedisCache{client: client}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := cache.Ping(ctx); err != nil {
		_ = cache.Close()
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return cache, nil
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, bool, error) {
	value, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return value, true, nil
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// IncrementWithExpiry atomically increments a rate-limit counter and ensures it has a TTL.
// Existing TTLs are preserved so the counter remains a fixed window.
func (r *RedisCache) IncrementWithExpiry(ctx context.Context, key string, window time.Duration) (int64, error) {
	ttlMilliseconds := window.Milliseconds()
	if ttlMilliseconds <= 0 {
		return 0, fmt.Errorf("rate-limit window must be at least 1ms, got %s", window)
	}

	count, err := incrementWithExpiryScript.Run(ctx, r.client, []string{key}, ttlMilliseconds).Int64()
	if err != nil {
		return 0, fmt.Errorf("execute atomic rate-limit counter script for key %q: %w", key, err)
	}
	return count, nil
}
