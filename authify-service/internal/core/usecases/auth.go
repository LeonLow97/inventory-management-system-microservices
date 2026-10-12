package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/LeonLow97/internal/core/domain"
	coreerrors "github.com/LeonLow97/internal/core/errors"
	"github.com/LeonLow97/internal/core/ports"
	"github.com/LeonLow97/internal/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Login performs authentication with bcrypt comparison of hashed password and provided password
func (s service) Login(ctx context.Context, loginInput domain.LoginInput) (*domain.User, string, error) {
	logger := s.logger.WithContext(ctx)
	user, err := s.repo.GetUserByEmail(ctx, loginInput.Email)
	if err != nil {
		if !errors.Is(err, coreerrors.ErrUserNotFound) {
			logger.WithError(err).Error("failed to look up user for login")
		}
		return nil, "", fmt.Errorf("look up user for login: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(utils.FromPointer(user.HashedPassword)),
		[]byte(loginInput.Password),
	); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) ||
			errors.Is(err, bcrypt.ErrHashTooShort):
			return nil, "", coreerrors.ErrInvalidCredentials
		default:
			logger.WithError(err).Error("password verification failed",
				ports.LogField{Key: "user_id", Value: user.ID},
			)
			return nil, "", fmt.Errorf("password verification failed for user %d: %w", user.ID, err)
		}
	}

	if !user.Active {
		return nil, "", coreerrors.ErrInactiveUser
	}

	// Create user session and store in cache
	sessionID := uuid.NewString()
	expiry := time.Duration(s.cfg.JWTConfig.Expiry) * time.Minute

	// Add small buffer to Redis TTL to handle clock drift and network latency
	const driftBuffer = time.Minute
	ttlSeconds := int64((expiry + driftBuffer).Seconds())

	if err := s.storeUserSessionInCache(ctx, sessionID, user.ID, ttlSeconds); err != nil {
		logger.WithError(err).Error("failed to store user session",
			ports.LogField{Key: "user_id", Value: user.ID},
		)
		return nil, "", err
	}

	token, err := s.generateJWTToken(user.ID, sessionID, expiry)
	if err != nil {
		logger.WithError(err).Error("failed to create authentication token",
			ports.LogField{Key: "user_id", Value: user.ID},
		)
		return nil, "", err
	}

	logger.Info("user authenticated", ports.LogField{Key: "user_id", Value: user.ID})
	return user, token, nil
}

func (s service) SignUp(ctx context.Context, signupInput domain.SignUpInput) error {
	logger := s.logger.WithContext(ctx)
	emailExists, err := s.repo.EmailExists(ctx, signupInput.Email)
	if err != nil {
		logger.WithError(err).Error("failed to check whether signup email exists")
		return fmt.Errorf("check whether signup email exists: %w", err)
	}
	if emailExists {
		return coreerrors.ErrEmailAlreadyExists
	}

	// Generate hashed password from plain text password via bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupInput.Password), 10)
	if err != nil {
		logger.WithError(err).Error("failed to hash signup password")
		return fmt.Errorf("hash signup password: %w", err)
	}

	user := &domain.User{
		Email:          signupInput.Email,
		HashedPassword: utils.ToPointer(string(hashedPassword)),
		FirstName:      utils.ToPointer(signupInput.FirstName),
		LastName:       utils.ToPointer(signupInput.LastName),
	}
	if err := s.repo.InsertUser(ctx, user); err != nil {
		logger.WithError(err).Error("failed to create user account")
		return fmt.Errorf("create user account: %w", err)
	}

	logger.Info("user account created")
	return nil
}

func (s service) generateJWTToken(userID int64, sessionID string, expiry time.Duration) (string, error) {
	expirationTime := time.Now().Add(expiry)
	claims := CustomClaims{
		SessionID: sessionID,
		UserID:    userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "authify-service",
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := generateToken.SignedString([]byte(s.cfg.JWTConfig.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token for user id '%d' with error: %v", userID, err)
	}
	return signedToken, nil
}

func (s service) storeUserSessionInCache(ctx context.Context, sessionID string, userID int64, ttlSeconds int64) error {
	key := fmt.Sprintf(userSessionCacheKey, sessionID)
	_, err := s.appCache.SetNX(ctx, key, userID, ttlSeconds)
	if err != nil {
		return fmt.Errorf("failed to set user session in cache with error: %w", err)
	}
	return nil
}
