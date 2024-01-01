package users

import "errors"

var (
	ErrNotFound              = errors.New("user not found")
	ErrSameValue             = errors.New("update of same value not allowed")
	ErrInvalidPasswordFormat = errors.New("invalid password format")
)
