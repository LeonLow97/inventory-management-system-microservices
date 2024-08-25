package ports

import (
	"context"

	"github.com/LeonLow97/internal/core/domain"
)

type UserRepo interface {
	GetUsers(ctx context.Context) (*[]domain.User, error)
	UpdateUser(ctx context.Context, req domain.User, userID int) error
}
