package order

import "errors"

var (
	ErrNoOrdersFound = errors.New("no orders found for this user")
)
