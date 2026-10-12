package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/LeonLow97/internal/infrastructure/config"
	"github.com/go-redis/redis/v8"
)

type redisCommanderStub struct {
	pingErr     error
	setErr      error
	setNXValue  bool
	setNXErr    error
	getValue    string
	getErr      error
	delErr      error
	existsCount int64
	existsErr   error
	lastKey     string
	lastKeys    []string
	lastValue   any
	lastTTL     time.Duration
}

func (s *redisCommanderStub) Ping(ctx context.Context) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx)
	cmd.SetVal("PONG")
	cmd.SetErr(s.pingErr)
	return cmd
}

func (s *redisCommanderStub) Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	s.lastKey, s.lastValue, s.lastTTL = key, value, expiration
	cmd := redis.NewStatusCmd(ctx)
	cmd.SetVal("OK")
	cmd.SetErr(s.setErr)
	return cmd
}

func (s *redisCommanderStub) SetNX(ctx context.Context, key string, value any, expiration time.Duration) *redis.BoolCmd {
	s.lastKey, s.lastValue, s.lastTTL = key, value, expiration
	cmd := redis.NewBoolCmd(ctx)
	cmd.SetVal(s.setNXValue)
	cmd.SetErr(s.setNXErr)
	return cmd
}

func (s *redisCommanderStub) Get(ctx context.Context, key string) *redis.StringCmd {
	s.lastKey = key
	cmd := redis.NewStringCmd(ctx)
	cmd.SetVal(s.getValue)
	cmd.SetErr(s.getErr)
	return cmd
}

func (s *redisCommanderStub) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	s.lastKeys = append([]string(nil), keys...)
	cmd := redis.NewIntCmd(ctx)
	cmd.SetVal(int64(len(keys)))
	cmd.SetErr(s.delErr)
	return cmd
}

func (s *redisCommanderStub) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	s.lastKeys = append([]string(nil), keys...)
	cmd := redis.NewIntCmd(ctx)
	cmd.SetVal(s.existsCount)
	cmd.SetErr(s.existsErr)
	return cmd
}

func (*redisCommanderStub) Close() error { return nil }

func TestRedisClientOptions(t *testing.T) {
	cfg := config.Config{RedisConfig: config.RedisConfig{
		Address: "redis.internal", Port: 6380, Password: "secret", DatabaseIndex: 4,
	}}
	options := redisClientOptions(cfg)
	if options.Addr != "redis.internal:6380" || options.Password != "secret" || options.DB != 4 {
		t.Errorf("Redis connection options = (%q, %q, %d)", options.Addr, options.Password, options.DB)
	}
	if options.MaxRetries != 3 || options.MinRetryBackoff != 100*time.Millisecond || options.MaxRetryBackoff != time.Second {
		t.Errorf("Redis retry options = (%d, %v, %v)", options.MaxRetries, options.MinRetryBackoff, options.MaxRetryBackoff)
	}
}

func TestRedisCachePing(t *testing.T) {
	testCases := []struct {
		name    string
		pingErr error
		wantErr bool
	}{
		{name: "success"},
		{name: "Redis error", pingErr: errors.New("connection failed"), wantErr: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cache := &RedisCache{client: &redisCommanderStub{pingErr: testCase.pingErr}}
			if err := cache.Ping(context.Background()); (err != nil) != testCase.wantErr {
				t.Errorf("Ping() error = %v, wantErr %t", err, testCase.wantErr)
			}
		})
	}
}

func TestRedisCacheSet(t *testing.T) {
	testCases := []struct {
		name    string
		setErr  error
		wantErr bool
	}{
		{name: "success"},
		{name: "Redis error", setErr: errors.New("write failed"), wantErr: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			stub := &redisCommanderStub{setErr: testCase.setErr}
			cache := &RedisCache{client: stub}
			err := cache.Set(context.Background(), "user:1", "value", 12)
			if (err != nil) != testCase.wantErr {
				t.Fatalf("Set() error = %v, wantErr %t", err, testCase.wantErr)
			}
			if stub.lastKey != "user:1" || stub.lastValue != "value" || stub.lastTTL != 12*time.Second {
				t.Errorf("Set() arguments = (%q, %v, %v), want (%q, %q, %v)", stub.lastKey, stub.lastValue, stub.lastTTL, "user:1", "value", 12*time.Second)
			}
		})
	}
}

func TestRedisCacheSetNX(t *testing.T) {
	testCases := []struct {
		name       string
		value      bool
		commandErr error
		wantValue  bool
		wantErr    bool
	}{
		{name: "key inserted", value: true, wantValue: true},
		{name: "key already exists"},
		{name: "Redis error", commandErr: errors.New("write failed"), wantErr: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			stub := &redisCommanderStub{setNXValue: testCase.value, setNXErr: testCase.commandErr}
			cache := &RedisCache{client: stub}
			got, err := cache.SetNX(context.Background(), "session:1", int64(42), 30)
			if got != testCase.wantValue || (err != nil) != testCase.wantErr {
				t.Fatalf("SetNX() = (%t, %v), want (%t, err=%t)", got, err, testCase.wantValue, testCase.wantErr)
			}
			if stub.lastKey != "session:1" || stub.lastValue != int64(42) || stub.lastTTL != 30*time.Second {
				t.Errorf("SetNX() arguments = (%q, %v, %v)", stub.lastKey, stub.lastValue, stub.lastTTL)
			}
		})
	}
}

func TestRedisCacheGet(t *testing.T) {
	readErr := errors.New("read failed")
	testCases := []struct {
		name       string
		value      string
		commandErr error
		wantValue  string
		wantErr    error
	}{
		{name: "returns value", value: "cached", wantValue: "cached"},
		{name: "returns missing key error", commandErr: redis.Nil, wantErr: redis.Nil},
		{name: "returns Redis error", commandErr: readErr, wantErr: readErr},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cache := &RedisCache{client: &redisCommanderStub{getValue: testCase.value, getErr: testCase.commandErr}}
			got, err := cache.Get(context.Background(), "session:1")
			if got != testCase.wantValue || (err != nil) != (testCase.wantErr != nil) {
				t.Fatalf("Get() = (%q, %v), want (%q, %v)", got, err, testCase.wantValue, testCase.wantErr)
			}
			if testCase.wantErr != nil && !errors.Is(err, testCase.wantErr) {
				t.Errorf("Get() error = %v, want %v", err, testCase.wantErr)
			}
		})
	}
}

func TestRedisCacheDelete(t *testing.T) {
	testCases := []struct {
		name    string
		delErr  error
		wantErr bool
	}{
		{name: "success"},
		{name: "Redis error", delErr: errors.New("delete failed"), wantErr: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			stub := &redisCommanderStub{delErr: testCase.delErr}
			cache := &RedisCache{client: stub}
			err := cache.Delete(context.Background(), "session:1")
			if (err != nil) != testCase.wantErr {
				t.Fatalf("Delete() error = %v, wantErr %t", err, testCase.wantErr)
			}
			if len(stub.lastKeys) != 1 || stub.lastKeys[0] != "session:1" {
				t.Errorf("Delete() keys = %v, want [session:1]", stub.lastKeys)
			}
		})
	}
}

func TestRedisCacheExists(t *testing.T) {
	testCases := []struct {
		name       string
		count      int64
		commandErr error
		want       bool
		wantErr    bool
	}{
		{name: "key exists", count: 1, want: true},
		{name: "key does not exist"},
		{name: "Redis error", commandErr: errors.New("lookup failed"), wantErr: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cache := &RedisCache{client: &redisCommanderStub{existsCount: testCase.count, existsErr: testCase.commandErr}}
			got, err := cache.Exists(context.Background(), "session:1")
			if got != testCase.want || (err != nil) != testCase.wantErr {
				t.Errorf("Exists() = (%t, %v), want (%t, err=%t)", got, err, testCase.want, testCase.wantErr)
			}
		})
	}
}
