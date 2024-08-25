package ports

import (
	"context"

	"github.com/LeonLow97/internal/core/domain"
)

type AuthRepo interface {
	Login(ctx context.Context, req domain.User) (*domain.User, error)
	SignUp(ctx context.Context, req domain.User) error
}
