package outbound

import "errors"

var (
	ErrOrdersNotFound = errors.New("orders not found")
	ErrOrderNotFound  = errors.New("order not found")

	ErrProductNotFound = errors.New("product not found for this user")
)
