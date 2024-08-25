package user

import (
	"context"
	"log"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/ports"
)

type User interface {
	GetUsers(ctx context.Context) (*[]domain.User, error)
	UpdateUser(ctx context.Context, req domain.User, userID int) error
}

type service struct {
	userRepo ports.UserRepo
}

func NewUserService(r ports.UserRepo) User {
	return &service{
		userRepo: r,
	}
}

func (s *service) GetUsers(ctx context.Context) (*[]domain.User, error) {
	users, err := s.userRepo.GetUsers(ctx)
	if err != nil {
		log.Printf("failed to get users with error: %v\n", err)
		return nil, err
	}

	return users, nil
}

func (s *service) UpdateUser(ctx context.Context, req domain.User, userID int) error {
	return s.userRepo.UpdateUser(ctx, req, userID)
}
