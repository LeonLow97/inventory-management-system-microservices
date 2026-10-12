package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"testing"
	"time"

	coreerrors "github.com/LeonLow97/internal/core/errors"
	"github.com/LeonLow97/internal/observability/logging"
)

const (
	testBurstKeyTemplate = "test:burst:%s:%s"
	testFailKeyTemplate  = "test:fail:%s:%s"
)

type memoryStore struct {
	values       map[string]string
	expiresAt    map[string]time.Time
	now          func() time.Time
	getErr       error
	deleteErr    error
	incrementErr error
}

func newMemoryStore() *memoryStore {
	return &memoryStore{
		values:    make(map[string]string),
		expiresAt: make(map[string]time.Time),
		now:       time.Now,
	}
}

func (s *memoryStore) Get(_ context.Context, key string) (string, bool, error) {
	if s.getErr != nil {
		return "", false, s.getErr
	}
	s.expireIfNeeded(key)
	value, ok := s.values[key]
	return value, ok, nil
}

func (s *memoryStore) Delete(_ context.Context, key string) error {
	if s.deleteErr != nil {
		return s.deleteErr
	}
	delete(s.values, key)
	delete(s.expiresAt, key)
	return nil
}

func (s *memoryStore) IncrementWithExpiry(_ context.Context, key string, window time.Duration) (int64, error) {
	if s.incrementErr != nil {
		return 0, s.incrementErr
	}
	if window < time.Millisecond {
		return 0, errors.New("rate-limit window must be at least 1ms")
	}
	if s.values == nil {
		s.values = make(map[string]string)
	}
	if s.expiresAt == nil {
		s.expiresAt = make(map[string]time.Time)
	}
	s.expireIfNeeded(key)

	count := int64(0)
	if value, found := s.values[key]; found {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, err
		}
		count = parsed
	}
	if _, hasExpiry := s.expiresAt[key]; !hasExpiry {
		s.expiresAt[key] = s.currentTime().Add(window)
	}
	count++
	s.values[key] = strconv.FormatInt(count, 10)
	return count, nil
}

func (s *memoryStore) expireIfNeeded(key string) {
	expiresAt, exists := s.expiresAt[key]
	if exists && !s.currentTime().Before(expiresAt) {
		delete(s.values, key)
		delete(s.expiresAt, key)
	}
}

func (s *memoryStore) currentTime() time.Time {
	if s.now == nil {
		return time.Now()
	}
	return s.now()
}

func newTestRateLimiter(store Store) *rateLimiter {
	cfg := AuthRateLimitingConfig{}
	cfg.Login.Burst = LoginRateLimitConfig{
		Key:         testBurstKeyTemplate,
		MaxRequests: 2,
		Window:      time.Minute,
	}
	cfg.Login.Fail = LoginRateLimitConfig{
		Key:         testFailKeyTemplate,
		MaxRequests: 3,
		Window:      5 * time.Minute,
	}
	return &rateLimiter{cfg: cfg, store: store, logger: logging.New(io.Discard)}
}

func TestFixedWindowCounter(t *testing.T) {
	now := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	store := newMemoryStore()
	store.now = func() time.Time { return now }
	limiter := newTestRateLimiter(store)

	for _, wantCount := range []int64{1, 2} {
		gotCount, err := limiter.fixedWindowCounter(context.Background(), "counter", 30*time.Second)
		if err != nil {
			t.Fatalf("fixedWindowCounter() error = %v", err)
		}
		if gotCount != wantCount {
			t.Errorf("fixedWindowCounter() count = %d, want %d", gotCount, wantCount)
		}
	}
	wantExpiry := now.Add(30 * time.Second)
	if got := store.expiresAt["counter"]; !got.Equal(wantExpiry) {
		t.Errorf("counter expiry = %v, want %v", got, wantExpiry)
	}

	now = wantExpiry
	gotCount, err := limiter.fixedWindowCounter(context.Background(), "counter", 30*time.Second)
	if err != nil {
		t.Fatalf("fixedWindowCounter() after window expiry error = %v", err)
	}
	if gotCount != 1 {
		t.Errorf("fixedWindowCounter() after window expiry count = %d, want 1", gotCount)
	}
	if got, want := store.expiresAt["counter"], now.Add(30*time.Second); !got.Equal(want) {
		t.Errorf("counter expiry after rollover = %v, want %v", got, want)
	}
}

func TestFixedWindowCounterErrors(t *testing.T) {
	testCases := []struct {
		name      string
		store     *memoryStore
		wantError string
	}{
		{name: "increment fails", store: &memoryStore{incrementErr: errors.New("redis increment failed")}, wantError: "failed to increment fixed-window counter \"counter\": redis increment failed"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			limiter := newTestRateLimiter(testCase.store)
			_, err := limiter.fixedWindowCounter(context.Background(), "counter", 30*time.Second)
			if err == nil || err.Error() != testCase.wantError {
				t.Fatalf("fixedWindowCounter() error = %v, want %q", err, testCase.wantError)
			}
		})
	}
}

func TestNewRateLimiterRejectsInvalidPolicy(t *testing.T) {
	validConfig := func() AuthRateLimitingConfig {
		return AuthRateLimitingConfig{
			Login: LoginRateLimitingConfig{
				Burst: LoginRateLimitConfig{Key: testBurstKeyTemplate, Window: time.Second, MaxRequests: 2},
				Fail:  LoginRateLimitConfig{Key: testFailKeyTemplate, Window: time.Minute, MaxRequests: 3},
			},
		}
	}

	testCases := []struct {
		name   string
		mutate func(*AuthRateLimitingConfig)
	}{
		{name: "zero burst maximum", mutate: func(cfg *AuthRateLimitingConfig) { cfg.Login.Burst.MaxRequests = 0 }},
		{name: "negative fail maximum", mutate: func(cfg *AuthRateLimitingConfig) { cfg.Login.Fail.MaxRequests = -1 }},
		{name: "zero window", mutate: func(cfg *AuthRateLimitingConfig) { cfg.Login.Burst.Window = 0 }},
		{name: "sub-millisecond window", mutate: func(cfg *AuthRateLimitingConfig) { cfg.Login.Fail.Window = time.Microsecond }},
		{name: "missing key template", mutate: func(cfg *AuthRateLimitingConfig) { cfg.Login.Burst.Key = "" }},
		{name: "invalid key template", mutate: func(cfg *AuthRateLimitingConfig) { cfg.Login.Fail.Key = "fail:%s" }},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := validConfig()
			testCase.mutate(&cfg)
			limiter, err := NewRateLimiter(cfg, newMemoryStore(), logging.New(io.Discard))
			if err == nil || limiter != nil {
				t.Fatalf("NewRateLimiter() = (%v, %v), want (nil, error)", limiter, err)
			}
		})
	}

	limiter, err := NewRateLimiter(validConfig(), newMemoryStore(), logging.New(io.Discard))
	if err != nil || limiter == nil {
		t.Fatalf("NewRateLimiter() = (%v, %v), want (limiter, nil)", limiter, err)
	}
}

func TestLoginRateLimiter(t *testing.T) {
	const email, ip = "ada@example.com", "192.0.2.1"
	burstKey := fmt.Sprintf(testBurstKeyTemplate, email, ip)
	failKey := fmt.Sprintf(testFailKeyTemplate, email, ip)

	testCases := []struct {
		name        string
		store       *memoryStore
		wantErr     bool
		wantErrKind coreerrors.Kind
		wantBurst   string
	}{
		{name: "allows request below limits", store: newMemoryStore(), wantBurst: "1"},
		{name: "rejects burst over limit", store: &memoryStore{values: map[string]string{burstKey: "2"}}, wantErrKind: coreerrors.KindResourceExhausted, wantBurst: "3"},
		{name: "rejects prior failed login limit", store: &memoryStore{values: map[string]string{failKey: "3"}}, wantErrKind: coreerrors.KindResourceExhausted, wantBurst: "1"},
		{name: "returns cache error when store is unavailable", store: &memoryStore{incrementErr: errors.New("redis unavailable"), getErr: errors.New("redis unavailable")}, wantErr: true},
		{name: "returns error for malformed fail counter", store: &memoryStore{values: map[string]string{failKey: "bad"}}, wantErr: true, wantBurst: "1"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			limiter := newTestRateLimiter(testCase.store)
			err := limiter.LoginRateLimiter(context.Background(), email, ip)
			if testCase.wantErr || testCase.wantErrKind != "" {
				if err == nil {
					t.Fatal("LoginRateLimiter() error = nil, want error")
				}
			}
			if testCase.wantErrKind != "" {
				var appErr *coreerrors.Error
				if !errors.As(err, &appErr) || appErr.Kind != testCase.wantErrKind {
					t.Errorf("LoginRateLimiter() error = %v, want kind %q", err, testCase.wantErrKind)
				}
			} else if !testCase.wantErr && err != nil {
				t.Errorf("LoginRateLimiter() error = %v, want nil", err)
			}
			if testCase.wantBurst != "" && testCase.store.values[burstKey] != testCase.wantBurst {
				t.Errorf("burst counter = %q, want %q", testCase.store.values[burstKey], testCase.wantBurst)
			}
		})
	}
}

func TestLoginFailRateLimitOperations(t *testing.T) {
	testCases := []struct {
		name              string
		prepare           func(*memoryStore, string)
		action            func(*rateLimiter, context.Context, string, string) error
		wantCounter       string
		wantCounterAbsent bool
		wantTTL           bool
	}{
		{
			name:        "increment creates counter and expiry",
			wantCounter: "1",
			wantTTL:     true,
			action: func(limiter *rateLimiter, ctx context.Context, email, ip string) error {
				return limiter.IncrementLoginFailRateLimiter(ctx, email, ip)
			},
		},
		{
			name:              "reset deletes counter",
			prepare:           func(store *memoryStore, key string) { store.values[key] = "2" },
			wantCounterAbsent: true,
			action: func(limiter *rateLimiter, ctx context.Context, email, ip string) error {
				return limiter.ResetLoginFailRateLimiter(ctx, email, ip)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			store := newMemoryStore()
			limiter := newTestRateLimiter(store)
			const email, ip = "ada@example.com", "192.0.2.1"
			key := fmt.Sprintf(testFailKeyTemplate, email, ip)
			if testCase.prepare != nil {
				testCase.prepare(store, key)
			}
			if err := testCase.action(limiter, context.Background(), email, ip); err != nil {
				t.Fatalf("operation error = %v", err)
			}
			gotCounter, exists := store.values[key]
			if testCase.wantCounterAbsent {
				if exists {
					t.Errorf("counter value = %q, want key deleted", gotCounter)
				}
			} else if gotCounter != testCase.wantCounter {
				t.Errorf("counter value = %q, want %q", gotCounter, testCase.wantCounter)
			}
			_, gotTTL := store.expiresAt[key]
			if gotTTL != testCase.wantTTL {
				t.Errorf("counter TTL exists = %t, want %t", gotTTL, testCase.wantTTL)
			}
		})
	}
}

func TestLoginFailRateLimitOperationsReturnStoreErrors(t *testing.T) {
	testCases := []struct {
		name   string
		store  *memoryStore
		action func(*rateLimiter) error
	}{
		{name: "increment", store: &memoryStore{incrementErr: errors.New("redis failed")}, action: func(l *rateLimiter) error { return l.IncrementLoginFailRateLimiter(context.Background(), "a", "b") }},
		{name: "reset", store: &memoryStore{deleteErr: errors.New("redis failed")}, action: func(l *rateLimiter) error { return l.ResetLoginFailRateLimiter(context.Background(), "a", "b") }},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if err := testCase.action(newTestRateLimiter(testCase.store)); err == nil {
				t.Fatal("expected store error")
			}
		})
	}
}
