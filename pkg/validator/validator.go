package validator

import (
	"github.com/go-playground/validator/v10"
)

// Validate is the shared validator instance.
var Validate = validator.New()

// FormatErrors extracts human-readable validation errors.
func FormatErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors[e.Field()] = formatTag(e)
		}
	}
	return errors
}

func formatTag(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email address"
	case "min":
		return "Value is too short (minimum: " + e.Param() + ")"
	case "max":
		return "Value is too long (maximum: " + e.Param() + ")"
	case "uuid":
		return "Must be a valid UUID"
	default:
		return "Invalid value"
	}
}
