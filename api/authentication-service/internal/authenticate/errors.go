package authenticate

import "errors"

var (
	ErrInvalidCredentials = errors.New("incorrect username and/or password")
	ErrInactiveUser       = errors.New("account is inactive")
	ErrNotFound           = errors.New("user not found")
)

var (
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrInvalidPasswordFormat = errors.New("invalid password format")
	ErrExistingUserFound = errors.New("existing user found")
)