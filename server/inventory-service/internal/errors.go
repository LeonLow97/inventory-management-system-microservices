package inventory

import "errors"

var (
	ErrProductNotFound = errors.New("product not found for this user")

	ErrBrandNotFound    = errors.New("brand not found")
	ErrCategoryNotFound = errors.New("category not found")
)
