package usecases

import (
	"context"
	"errors"

	"github.com/LeonLow97/internal/core/domain"
	coreerrors "github.com/LeonLow97/internal/core/errors"
	"github.com/LeonLow97/internal/core/ports"
)

type User interface {
	GetCurrentUser(ctx context.Context) (*domain.User, error)
	Login(ctx context.Context, req domain.User, clientIP string) (*domain.User, string, error)
	SignUp(ctx context.Context, req domain.User) error

	GetUsers(ctx context.Context, limit int64, cursor string) ([]domain.User, string, error)
	UpdateUser(ctx context.Context, req domain.User) error
}

type service struct {
	userRepo      ports.UserOutbound
	loginAttempts ports.LoginAttemptLimiter
	logger        ports.Logger
}

func NewUserService(userRepo ports.UserOutbound, loginAttempts ports.LoginAttemptLimiter, logger ports.Logger) User {
	return &service{
		userRepo:      userRepo,
		loginAttempts: loginAttempts,
		logger:        logger,
	}
}

func (s *service) GetCurrentUser(ctx context.Context) (*domain.User, error) {
	return s.userRepo.GetCurrentUser(ctx)
}

func (s *service) Login(ctx context.Context, req domain.User, clientIP string) (*domain.User, string, error) {
	logger := s.logger.WithContext(ctx)
	if err := s.loginAttempts.LoginRateLimiter(ctx, req.Email, clientIP); err != nil {
		if isRateLimitExceeded(err) {
			return nil, "", err
		}
		logger.WithError(err).Warn("login rate-limit check failed; continuing")
	}

	user, token, err := s.userRepo.Login(ctx, req)
	if err != nil {
		if isUnauthenticated(err) {
			if counterErr := s.loginAttempts.IncrementLoginFailRateLimiter(ctx, req.Email, clientIP); counterErr != nil {
				logger.WithError(counterErr).Warn("failed to increment login failure counter")
			}
		}
		logger.WithError(err).Warn("login failed",
			ports.LogField{Key: "email", Value: req.Email},
		)
		return nil, "", err
	}

	if err := s.loginAttempts.ResetLoginFailRateLimiter(ctx, req.Email, clientIP); err != nil {
		logger.WithError(err).Warn("failed to reset login failure counter")
	}

	return user, token, nil
}

func isRateLimitExceeded(err error) bool {
	var appErr *coreerrors.Error
	return errors.As(err, &appErr) && appErr.Kind == coreerrors.KindResourceExhausted
}

func isUnauthenticated(err error) bool {
	var appErr *coreerrors.Error
	return errors.As(err, &appErr) && appErr.Kind == coreerrors.KindUnauthenticated
}

func (s *service) SignUp(ctx context.Context, req domain.User) error {
	logger := s.logger.WithContext(ctx)
	if err := s.userRepo.SignUp(ctx, req); err != nil {
		logger.WithError(err).Warn("sign-up failed",
			ports.LogField{Key: "email", Value: req.Email},
		)
		return err
	}
	return nil
}

func (s *service) GetUsers(ctx context.Context, limit int64, cursor string) ([]domain.User, string, error) {
	logger := s.logger.WithContext(ctx)
	users, nextCursor, err := s.userRepo.GetUsers(ctx, limit, cursor)
	if err != nil {
		logger.WithError(err).Error("failed to get users")
		return nil, "", err
	}
	return users, nextCursor, nil
}

func (s *service) UpdateUser(ctx context.Context, req domain.User) error {
	return s.userRepo.UpdateUser(ctx, req)
}
