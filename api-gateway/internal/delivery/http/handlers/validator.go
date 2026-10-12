package handlers

import (
	"strings"

	"github.com/LeonLow97/internal/core/ports"
	"github.com/LeonLow97/internal/delivery/http/validation"
	"github.com/go-playground/validator/v10"
)

func (h BaseHandler) ValidateStruct(v any) error {
	return h.validator.Struct(v)
}

// validatorInstance creates and returns a new validator instance with custom validations
func validatorInstance(logger ports.Logger) *validator.Validate {
	validatorInstance := validator.New()

	// Valid Password format
	if err := validatorInstance.RegisterValidation("password_format", func(f1 validator.FieldLevel) bool {
		return validation.IsValidPassword(strings.TrimSpace(f1.Field().String()))
	}); err != nil {
		logger.WithError(err).Error("failed to register password format validation")
	}

	return validatorInstance
}
