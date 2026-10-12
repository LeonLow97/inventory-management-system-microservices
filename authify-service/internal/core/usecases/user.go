package usecases

import (
	"context"
	"fmt"

	"github.com/LeonLow97/internal/core/domain"
	coreerrors "github.com/LeonLow97/internal/core/errors"
	"github.com/LeonLow97/internal/core/ports"
	"github.com/LeonLow97/internal/pkg/cursor"
	"golang.org/x/crypto/bcrypt"
)

func (s service) GetUsers(ctx context.Context, adminUserID int64, limit int64, cursorStr string) ([]domain.User, string, error) {
	logger := s.logger.WithContext(ctx)
	// Check if user is an admin user
	isAdminUser, err := s.repo.IsAdminUser(ctx, adminUserID)
	if err != nil {
		logger.WithError(err).Error("failed to check admin access",
			ports.LogField{Key: "user_id", Value: adminUserID},
		)
		return nil, "", err
	}
	if !isAdminUser {
		logger.Warn("admin-only user listing denied", ports.LogField{Key: "user_id", Value: adminUserID})
		return nil, "", coreerrors.ErrUnauthorizedAdminAccess
	}

	// Decode the cursor
	var userCursor domain.UserCursor
	if err := cursor.DecodeCursor(cursorStr, &userCursor); err != nil {
		logger.WithError(err).Warn("failed to decode user-list cursor")
		return nil, "", err
	}

	users, err := s.repo.GetUsers(ctx, limit, userCursor)
	if err != nil {
		logger.WithError(err).Error("failed to list users",
			ports.LogField{Key: "user_id", Value: adminUserID},
		)
		return nil, "", err
	}

	// Encode the next cursor based on the last user's ID
	var nextCursor string
	if len(users) > 0 {
		lastUserID := users[len(users)-1].ID
		nextCursor, err = cursor.EncodeCursor(domain.UserCursor{ID: lastUserID})
		if err != nil {
			logger.WithError(err).Error("failed to encode next user-list cursor")
			return nil, "", err
		}
	}

	logger.Info("users listed",
		ports.LogField{Key: "admin_user_id", Value: adminUserID},
		ports.LogField{Key: "user_count", Value: len(users)},
	)
	return users, nextCursor, nil
}

func (s service) UpdateUserByEmail(ctx context.Context, userID int64, updateUserInput domain.UpdateUserInput) error {
	logger := s.logger.WithContext(ctx)
	var hashedPassword *string
	if updateUserInput.Password != nil && *updateUserInput.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*updateUserInput.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.WithError(err).Error("failed to hash updated password",
				ports.LogField{Key: "user_id", Value: userID},
			)
			return fmt.Errorf("hash updated password: %w", err)
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

	if err := s.repo.UpdateUserByEmail(ctx, updateUser); err != nil {
		logger.WithError(err).Error("failed to update user",
			ports.LogField{Key: "user_id", Value: userID},
		)
		return fmt.Errorf("update user: %w", err)
	}

	logger.Info("user updated", ports.LogField{Key: "user_id", Value: userID})
	return nil
}
