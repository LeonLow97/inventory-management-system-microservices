package ports

import (
	"context"

	"github.com/LeonLow97/internal/core/domain"
)

type UserRepo interface {
	GetUsers(ctx context.Context, limit int64, cursor string) ([]domain.User, string, error)
	UpdateUser(ctx context.Context, req domain.User) error
}
