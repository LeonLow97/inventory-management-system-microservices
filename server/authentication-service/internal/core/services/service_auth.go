package services

import (
	"context"
	"log"
	"time"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func (s service) Login(ctx context.Context, username, password string) (*domain.User, string, error) {
	// retrieve user details by username
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		log.Printf("failed to retrieve user by username '%s' during login with error: %v\n", username, err)
		return nil, "", err
	}

	// compare db hashed password with user supplied plain text password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		switch {
		case err == bcrypt.ErrMismatchedHashAndPassword,
			err == bcrypt.ErrHashTooShort:
			log.Printf("invalid credentials for user id '%d' when performing bcrypt comparison\n", user.ID)
			return nil, "", ErrInvalidCredentials
		default:
			log.Printf("failed to compare hash and password with error: %v\n", err)
			return nil, "", err
		}
	}

	// check user account status
	if user.Active == 0 {
		log.Printf("user id '%d' is inactive\n", user.ID)
		return nil, "", ErrInactiveUser
	}

	// generate JWT token for user
	token, err := utils.GenerateJWTToken(user, time.Duration(s.cfg.JWTConfig.Expiry)*time.Minute, s.cfg.JWTConfig.SecretKey)
	if err != nil {
		log.Printf("failed to generate JWT token for user id '%d'\n", user.ID)
		return nil, "", err
	}

	return user, token, err
}

func (s service) SignUp(ctx context.Context, user *domain.User) error {
	if !utils.IsValidPassword(user.Email) {
		log.Printf("invalid email format '%s'\n", user.Email)
		return ErrInvalidEmailFormat
	}
	if !utils.IsValidPassword(user.Password) {
		log.Println("invalid password format")
		return ErrInvalidPasswordFormat
	}

	taken, err := s.repo.IsUsernameTaken(ctx, user.Username)
	if err != nil {
		log.Printf("failed to check if username has been taken for username '%s' with error: %v\n", user.Username, err)
		return err
	}

	if taken {
		return ErrUsernameTaken
	}

	// bcrypt hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Printf("failed to hash password for username '%s' with error: %v\n", user.Username, err)
		return err
	}
	user.Password = string(hashedPassword)

	// create user
	if err := s.repo.InsertUser(ctx, user); err != nil {
		return err
	}

	return nil
}
