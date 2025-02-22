package services

import (
	"context"
	"log"
	"time"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// Login performs authentication with bcrypt comparison of hashed password and provided password
func (s service) Login(ctx context.Context, loginInput domain.LoginInput) (*domain.User, string, error) {
	// Retrieve user details by username
	user, err := s.repo.GetUserByUsername(ctx, loginInput.Username)
	if err != nil {
		log.Printf("failed to retrieve user by username '%s' during login with error: %v\n", loginInput.Username, err)
		return nil, "", err
	}

	// Compare hashed password with user supplied plain text password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password)); err != nil {
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

	// Check if user is inactive
	if user.Active == 0 {
		log.Printf("user id '%d' is inactive\n", user.ID)
		return nil, "", ErrInactiveUser
	}

	// Generate JWT token for user
	token, err := utils.GenerateJWTToken(user, time.Duration(s.cfg.JWTConfig.Expiry)*time.Minute, s.cfg.JWTConfig.SecretKey)
	if err != nil {
		log.Printf("failed to generate JWT token for user id '%d'\n", user.ID)
		return nil, "", err
	}

	return user, token, err
}

func (s service) SignUp(ctx context.Context, signupInput domain.SignUpInput) error {
	// TODO: Shift this to the API Gateway to perform validation
	if !utils.IsValidPassword(signupInput.Email) {
		log.Printf("invalid email format '%s'\n", signupInput.Email)
		return ErrInvalidEmailFormat
	}
	// TODO: Shift this to the API Gateway to perform validation
	if !utils.IsValidPassword(signupInput.Password) {
		log.Println("invalid password format")
		return ErrInvalidPasswordFormat
	}

	taken, err := s.repo.IsUsernameTaken(ctx, signupInput.Username)
	if err != nil {
		log.Printf("failed to check if username has been taken for username '%s' with error: %v\n", signupInput.Username, err)
		return err
	}
	if taken {
		return ErrUsernameTaken
	}

	// bcrypt hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupInput.Password), 10)
	if err != nil {
		log.Printf("failed to hash password for username '%s' with error: %v\n", signupInput.Username, err)
		return err
	}

	// create user
	user := &domain.User{
		FirstName: signupInput.FirstName,
		LastName:  signupInput.LastName,
		Username:  signupInput.Username,
		Password:  string(hashedPassword),
		Email:     signupInput.Email,
	}
	if err := s.repo.InsertUser(ctx, user); err != nil {
		log.Printf("failed to insert user during signup with error: %v\n", err)
		return err
	}

	return nil
}
