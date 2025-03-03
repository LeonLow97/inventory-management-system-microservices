package handler

import "github.com/go-playground/validator/v10"

type Handler struct {
	validator *validator.Validate
}

func NewHandler() Handler {
	return Handler{
		validator: validatorInstance(),
	}
}
