package services

import (
	"context"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/pkg/config"
	"github.com/LeonLow97/internal/ports"
)

type Service interface {
	Login(ctx context.Context, username, password string) (*domain.User, string, error)
	SignUp(ctx context.Context, user *domain.User) error

	UpdateUser(ctx context.Context, user *domain.User) error
	GetUsers(ctx context.Context) (*[]domain.User, error)
}

type service struct {
	repo ports.Repository
	cfg  config.Config
}

func NewService(r ports.Repository, cfg config.Config) Service {
	return &service{
		repo: r,
		cfg:  cfg,
	}
}
