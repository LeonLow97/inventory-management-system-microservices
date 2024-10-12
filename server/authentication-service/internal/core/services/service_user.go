package services

import (
	"context"
	"log"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func (s service) UpdateUser(ctx context.Context, user *domain.User) error {
	user, err := s.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		log.Printf("failed to retrieve user by id '%d' with error: %v\n", user.ID, err)
		return err
	}

	// check if password is valid (when provided)
	if len(user.Password) > 0 && !utils.IsValidPassword(user.Password) {
		return ErrInvalidPasswordFormat
	} else {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("failed to generate updated password for user id '%d' with error: %v\n", user.ID, err)
			return err
		}
		user.Password = string(hashedPassword)
	}

	if err := s.repo.UpdateUserByID(ctx, user); err != nil {
		log.Printf("failed to update user by id '%d' with error: %v\n", user.ID, err)
		return err
	}

	return nil
}

// TODO: use cursor pagination for this
func (s service) GetUsers(ctx context.Context) (*[]domain.User, error) {
	return &[]domain.User{}, nil
}
