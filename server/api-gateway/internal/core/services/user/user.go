package user

import (
	"context"
	"log"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/ports"
)

type User interface {
	GetUsers(ctx context.Context, limit int64, cursor string) ([]domain.User, string, error)
	UpdateUser(ctx context.Context, req domain.User) error
}

type service struct {
	userRepo ports.UserRepo
}

func NewUserService(r ports.UserRepo) User {
	return &service{
		userRepo: r,
	}
}

func (s *service) GetUsers(ctx context.Context, limit int64, cursor string) ([]domain.User, string, error) {
	users, nextCursor, err := s.userRepo.GetUsers(ctx, limit, cursor)
	if err != nil {
		log.Printf("failed to get users with error: %v\n", err)
		return nil, "", err
	}

	return users, nextCursor, nil
}

func (s *service) UpdateUser(ctx context.Context, req domain.User) error {
	return s.userRepo.UpdateUser(ctx, req)
}
