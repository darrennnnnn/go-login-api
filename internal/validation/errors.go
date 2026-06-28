package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormatValidationErrors(err error) []ValidationError {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return []ValidationError{
			{
				Field:   "request",
				Message: "Invalid request body.",
			},
		}
	}

	errors := make([]ValidationError, 0, len(validationErrors))

	for _, e := range validationErrors {
		field := strings.ToLower(e.Field())

		var message string

		switch e.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required.", field)

		case "min":
			message = fmt.Sprintf("%s must be at least %s characters long.", field, e.Param())

		case "max":
			message = fmt.Sprintf("%s must be at most %s characters long.", field, e.Param())

		case "email":
			message = fmt.Sprintf("%s must be a valid email address.", field)

		case "password":
			message = "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character."

		default:
			message = fmt.Sprintf("%s is invalid.", field)
		}

		errors = append(errors, ValidationError{
			Field:   field,
			Message: message,
		})
	}

	return errors
}