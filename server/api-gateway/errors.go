package main

import "errors"

var (
	ErrMissingUserIDInJWTToken = errors.New("missing user id in jwt token")
)
