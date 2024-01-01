package authenticate

import (
	"fmt"
	"log"
	"time"

	"github.com/LeonLow97/utils"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(req LoginRequestDTO) (*User, string, error)
	SignUp(req SignUpRequestDTO) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s service) Login(req LoginRequestDTO) (*User, string, error) {
	// check if username exists
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, "nil", err
	}

	// compare password with hashed password in db
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		switch {
		case err == bcrypt.ErrMismatchedHashAndPassword, err == bcrypt.ErrHashTooShort:
			return nil, "", ErrInvalidCredentials
		default:
			return nil, "", err
		}
	}

	// check account status (active / inactive)
	if user.Active == 0 {
		return nil, "", ErrInactiveUser
	}

	// generate jwt token for user
	token, err := generateJWTToken(user)
	if err != nil {
		return &user, "", err
	}

	return &user, token, nil
}

func (s service) SignUp(req SignUpRequestDTO) error {
	// check format of email address
	if !utils.IsValidEmail(req.Email) {
		return ErrInvalidEmailFormat
	}

	// check if username already exists in database
	count, err := s.repo.GetUserCountByUsername(req.Username)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrExistingUserFound
	}

	// check if password contains at least 1 uppercase, lowercase, numeric and special character
	if !utils.IsValidPassword(req.Password) {
		return ErrInvalidPasswordFormat
	}

	// hash password with salt rounds of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return err
	}
	req.Password = string(hashedPassword)

	// insert into users table
	if err := s.repo.InsertOneUser(req); err != nil {
		return err
	}

	return nil
}

func generateJWTToken(user User) (string, error) {
	// generate token with claims
	tokenExpireTime := time.Now().Add(1 * time.Hour)
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    fmt.Sprintf("%d", user.ID),
		ExpiresAt: jwt.NewNumericDate(tokenExpireTime), // 1 hour
	})

	signedToken, err := generateToken.SignedString([]byte(JWT_SECRET_KEY))
	if err != nil {
		log.Println("error generating jwt token", err)
		return "", err
	}

	return signedToken, nil
}
