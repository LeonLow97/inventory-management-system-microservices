package inventory

import "errors"

var (
	ErrProductsNotFound = errors.New("no products found for the given userID")
	ErrProductNotFound  = errors.New("product not found for this user")

	ErrBrandNotFound    = errors.New("brand not found")
	ErrCategoryNotFound = errors.New("category not found")
)
