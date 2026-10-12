package usecases

import (
	"context"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/core/ports"
)

type Service interface {
	Login(ctx context.Context, loginInput domain.LoginInput) (*domain.User, string, error)
	SignUp(ctx context.Context, signupInput domain.SignUpInput) error

	UpdateUserByEmail(ctx context.Context, userID int64, updateUserInput domain.UpdateUserInput) error
	GetUsers(ctx context.Context, adminUserID int64, limit int64, cursorStr string) ([]domain.User, string, error)
}

type service struct {
	repo     ports.UserRepository
	cfg      Config
	appCache ports.Cache
	logger   ports.Logger
}

type Config struct {
	JWTConfig JWTConfig
}

type JWTConfig struct {
	SecretKey string
	Expiry    int
}

func NewService(cfg Config, r ports.UserRepository, appCache ports.Cache, logger ports.Logger) Service {
	return &service{
		repo:     r,
		cfg:      cfg,
		appCache: appCache,
		logger:   logger,
	}
}
