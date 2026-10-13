package ports

import (
	"context"

	"github.com/LeonLow97/internal/core/domain"
)

type UserOutbound interface {
	GetCurrentUser(ctx context.Context) (*domain.User, error)
	Login(ctx context.Context, req domain.User) (*domain.User, string, error)
	SignUp(ctx context.Context, req domain.User) error

	GetUsers(ctx context.Context, limit int64, cursor string) ([]domain.User, string, error)
	UpdateUser(ctx context.Context, req domain.User) error
}

// LoginAttemptLimiter checks login policy and updates attempt counters.
type LoginAttemptLimiter interface {
	LoginRateLimiter(ctx context.Context, email, ip string) error
	ResetLoginFailRateLimiter(ctx context.Context, email, ip string) error
	IncrementLoginFailRateLimiter(ctx context.Context, email, ip string) error
}
