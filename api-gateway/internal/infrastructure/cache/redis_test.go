package cache

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/LeonLow97/internal/infrastructure/config"
	"github.com/go-redis/redis/v8"
)

type redisCommanderStub struct {
	pingErr      error
	getValue     string
	getErr       error
	delErr       error
	evalValue    int64
	evalErr      error
	closeErr     error
	lastKey      string
	lastKeys     []string
	lastArgs     []interface{}
	closeCalls   int
	pingDeadline bool
}

func statusCommand(ctx context.Context, value string, err error) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx)
	cmd.SetVal(value)
	cmd.SetErr(err)
	return cmd
}

func stringCommand(ctx context.Context, value string, err error) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx)
	cmd.SetVal(value)
	cmd.SetErr(err)
	return cmd
}

func intCommand(ctx context.Context, value int64, err error) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)
	cmd.SetVal(value)
	cmd.SetErr(err)
	return cmd
}

func (s *redisCommanderStub) Ping(ctx context.Context) *redis.StatusCmd {
	_, s.pingDeadline = ctx.Deadline()
	return statusCommand(ctx, "PONG", s.pingErr)
}

func (s *redisCommanderStub) Get(ctx context.Context, key string) *redis.StringCmd {
	s.lastKey = key
	return stringCommand(ctx, s.getValue, s.getErr)
}

func (s *redisCommanderStub) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	s.lastKeys = append([]string(nil), keys...)
	return intCommand(ctx, 1, s.delErr)
}

func (s *redisCommanderStub) Eval(ctx context.Context, _ string, keys []string, args ...interface{}) *redis.Cmd {
	s.lastKeys = append([]string(nil), keys...)
	s.lastArgs = append([]interface{}(nil), args...)
	command := redis.NewCmd(ctx)
	command.SetVal(s.evalValue)
	command.SetErr(s.evalErr)
	return command
}

func (s *redisCommanderStub) EvalSha(ctx context.Context, _ string, _ []string, _ ...interface{}) *redis.Cmd {
	command := redis.NewCmd(ctx)
	command.SetErr(errors.New("NOSCRIPT script is not loaded"))
	return command
}

func (s *redisCommanderStub) ScriptExists(ctx context.Context, _ ...string) *redis.BoolSliceCmd {
	command := redis.NewBoolSliceCmd(ctx)
	command.SetVal([]bool{false})
	return command
}

func (s *redisCommanderStub) ScriptLoad(ctx context.Context, _ string) *redis.StringCmd {
	command := redis.NewStringCmd(ctx)
	command.SetVal("script-sha")
	return command
}

func (s *redisCommanderStub) Close() error {
	s.closeCalls++
	return s.closeErr
}

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

func TestNewRedisCachePingsAndClosesOnFailure(t *testing.T) {
	testCases := []struct {
		name    string
		pingErr error
		wantErr bool
	}{
		{name: "ping succeeds"},
		{name: "ping failure closes client", pingErr: errors.New("connection refused"), wantErr: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			stub := &redisCommanderStub{pingErr: testCase.pingErr}
			cache, err := newRedisCache(stub)
			if (err != nil) != testCase.wantErr {
				t.Fatalf("newRedisCache() error = %v, wantErr %t", err, testCase.wantErr)
			}
			if !stub.pingDeadline {
				t.Error("startup ping did not use a bounded context")
			}
			if testCase.wantErr {
				if cache != nil || stub.closeCalls != 1 {
					t.Errorf("failed initialization returned cache=%v and closed client %d times", cache, stub.closeCalls)
				}
				if !errors.Is(err, testCase.pingErr) {
					t.Errorf("initialization error = %v, want to wrap %v", err, testCase.pingErr)
				}
			} else if cache == nil || stub.closeCalls != 0 {
				t.Errorf("successful initialization returned cache=%v and closed client %d times", cache, stub.closeCalls)
			}
		})
	}
}

func TestRedisCachePingAndClose(t *testing.T) {
	testCases := []struct {
		name         string
		pingErr      error
		closeErr     error
		wantPingErr  bool
		wantCloseErr bool
	}{
		{name: "successful connection", wantPingErr: false, wantCloseErr: false},
		{name: "ping error", pingErr: errors.New("redis unavailable"), wantPingErr: true},
		{name: "close error", closeErr: errors.New("close failed"), wantCloseErr: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			stub := &redisCommanderStub{pingErr: testCase.pingErr, closeErr: testCase.closeErr}
			cache := &RedisCache{client: stub}
			if gotErr := cache.Ping(context.Background()); (gotErr != nil) != testCase.wantPingErr {
				t.Errorf("Ping() error = %v, wantErr %t", gotErr, testCase.wantPingErr)
			}
			if gotErr := cache.Close(); (gotErr != nil) != testCase.wantCloseErr {
				t.Errorf("Close() error = %v, wantErr %t", gotErr, testCase.wantCloseErr)
			}
			if stub.closeCalls != 1 {
				t.Errorf("Close calls = %d, want 1", stub.closeCalls)
			}
		})
	}
}

func TestRedisCacheGet(t *testing.T) {
	testCases := []struct {
		name       string
		value      string
		commandErr error
		wantValue  string
		wantFound  bool
		wantErr    bool
	}{
		{name: "found", value: "cached", wantValue: "cached", wantFound: true},
		{name: "missing key", commandErr: redis.Nil, wantValue: "", wantFound: false},
		{name: "redis failure", commandErr: errors.New("connection lost"), wantErr: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cache := &RedisCache{client: &redisCommanderStub{getValue: testCase.value, getErr: testCase.commandErr}}
			gotValue, gotFound, err := cache.Get(context.Background(), "cache-key")
			if gotValue != testCase.wantValue || gotFound != testCase.wantFound || (err != nil) != testCase.wantErr {
				t.Errorf("Get() = (%q, %t, %v), want (%q, %t, err=%t)", gotValue, gotFound, err, testCase.wantValue, testCase.wantFound, testCase.wantErr)
			}
		})
	}
}

func TestRedisCacheIncrementWithExpiry(t *testing.T) {
	ctx := context.Background()
	stub := &redisCommanderStub{evalValue: 4}
	cache := &RedisCache{client: stub}

	got, err := cache.IncrementWithExpiry(ctx, "counter", 1500*time.Millisecond)
	if err != nil || got != 4 {
		t.Fatalf("IncrementWithExpiry() = (%d, %v), want (4, nil)", got, err)
	}
	if !reflect.DeepEqual(stub.lastKeys, []string{"counter"}) {
		t.Errorf("script keys = %v, want [counter]", stub.lastKeys)
	}
	if !reflect.DeepEqual(stub.lastArgs, []interface{}{int64(1500)}) {
		t.Errorf("script arguments = %v, want [1500 milliseconds]", stub.lastArgs)
	}
}

func TestRedisCacheIncrementWithExpiryRejectsSubMillisecondWindow(t *testing.T) {
	stub := &redisCommanderStub{}
	cache := &RedisCache{client: stub}
	if _, err := cache.IncrementWithExpiry(context.Background(), "counter", time.Microsecond); err == nil {
		t.Fatal("IncrementWithExpiry() error = nil, want sub-millisecond window error")
	}
	if len(stub.lastKeys) != 0 {
		t.Errorf("script keys = %v, want no script call", stub.lastKeys)
	}
}

func TestRedisCacheIncrementWithExpiryPropagatesScriptError(t *testing.T) {
	stub := &redisCommanderStub{evalErr: errors.New("script failed")}
	cache := &RedisCache{client: stub}
	if _, err := cache.IncrementWithExpiry(context.Background(), "counter", time.Second); err == nil {
		t.Fatal("IncrementWithExpiry() error = nil, want script error")
	}
}

func TestRedisCacheDelete(t *testing.T) {
	ctx := context.Background()
	t.Run("Delete passes keys", func(t *testing.T) {
		stub := &redisCommanderStub{}
		if err := (&RedisCache{client: stub}).Delete(ctx, "one"); err != nil {
			t.Fatalf("Delete() error = %v", err)
		}
		if !reflect.DeepEqual(stub.lastKeys, []string{"one"}) {
			t.Errorf("Delete keys = %v, want [one]", stub.lastKeys)
		}
	})
	testCasesWithErrors := []struct {
		name   string
		cache  *RedisCache
		action func(*RedisCache) error
	}{
		{
			name:   "Delete propagates error",
			cache:  &RedisCache{client: &redisCommanderStub{delErr: errors.New("delete failed")}},
			action: func(cache *RedisCache) error { return cache.Delete(ctx, "key") },
		},
	}
	for _, testCase := range testCasesWithErrors {
		t.Run(testCase.name, func(t *testing.T) {
			if err := testCase.action(testCase.cache); err == nil {
				t.Fatal("expected Redis command error")
			}
		})
	}
}
