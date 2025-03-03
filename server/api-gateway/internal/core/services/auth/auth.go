package auth

import (
	"context"
	"log"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/ports"
)

type Auth interface {
	Login(ctx context.Context, req domain.User) (*domain.User, error)
	SignUp(ctx context.Context, req domain.User) error
}

type service struct {
	authRepo ports.AuthRepo
}

func NewAuthService(r ports.AuthRepo) Auth {
	return &service{
		authRepo: r,
	}
}

func (s *service) Login(ctx context.Context, req domain.User) (*domain.User, error) {
	user, err := s.authRepo.Login(ctx, req)
	if err != nil {
		log.Printf("failed to login for email %s with error: %v\n", req.Email, err)
		return nil, err
	}
	return user, nil
}

func (s *service) SignUp(ctx context.Context, req domain.User) error {
	if err := s.authRepo.SignUp(ctx, req); err != nil {
		log.Printf("failed to sign up for email %s with error: %v\n", req.Email, err)
		return err
	}
	return nil
}
