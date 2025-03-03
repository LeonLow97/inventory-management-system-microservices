package services

import (
	"context"
	"log"

	"github.com/LeonLow97/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
)

func (s service) UpdateUser(ctx context.Context, userID int64, updateUserInput domain.UpdateUserInput) error {
	var hashedPassword *string
	if updateUserInput.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*updateUserInput.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("failed to generate updated password for user id '%d' with error: %v\n", userID, err)
			return err
		}
		hashedStr := string(hashed)
		hashedPassword = &hashedStr
	}

	updateUser := &domain.User{
		ID:             userID,
		HashedPassword: hashedPassword,
		FirstName:      updateUserInput.FirstName,
		LastName:       updateUserInput.LastName,
	}

	if err := s.repo.UpdateUserByID(ctx, updateUser); err != nil {
		log.Printf("failed to update user by id '%d' with error: %v\n", userID, err)
		return err
	}

	return nil
}

// TODO: use cursor pagination for this
func (s service) GetUsers(ctx context.Context) (*[]domain.User, error) {
	return &[]domain.User{}, nil
}
