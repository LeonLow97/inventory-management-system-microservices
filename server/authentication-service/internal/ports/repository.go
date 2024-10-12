package ports

import (
	"context"

	"github.com/LeonLow97/internal/core/domain"
)

type Repository interface {
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByID(ctx context.Context, userID int) (*domain.User, error)
	GetUsers() (*[]domain.User, error)
	IsUsernameTaken(ctx context.Context, username string) (bool, error)

	InsertUser(ctx context.Context, user *domain.User) error
	UpdateUserByID(ctx context.Context, user *domain.User) error
}
