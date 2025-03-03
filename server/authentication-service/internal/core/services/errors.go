package services

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInactiveUser       = errors.New("inactive user")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")

	ErrInvalidEmailFormat    = errors.New("invalid email format")
	ErrInvalidPasswordFormat = errors.New("invalid password format")
)
