package inventory

import "errors"

var (
	ErrProductNotFound = errors.New("product not found for this user")
	
	ErrBrandOrCategoryNotFound = errors.New("brand or category not found")
)