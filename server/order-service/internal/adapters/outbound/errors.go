package outbound

import "errors"

var (
	ErrNoOrdersFound = errors.New("no orders found for this user")
	ErrNoOrderFound  = errors.New("unable to find the order for this user")

	ErrProductNotFound = errors.New("product not found for this user")
)
