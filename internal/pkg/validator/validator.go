package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Validate validates a struct and returns formatted field errors
func Validate(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errs := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		errs[field] = formatError(e)
	}
	return errs
}

func formatError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return "invalid email address"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", e.Field(), e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s failed validation: %s", e.Field(), e.Tag())
	}
}
