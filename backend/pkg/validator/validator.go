package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var defaultValidate = validator.New()

type EchoValidator struct {
	validate *validator.Validate
}

func New() *EchoValidator {
	return &EchoValidator{validate: validator.New()}
}

func (v *EchoValidator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func ValidateStruct(i interface{}) error {
	if i == nil {
		return fmt.Errorf("value is required")
	}
	return defaultValidate.Struct(i)
}

func FormatValidationError(err error) error {
	if err == nil {
		return nil
	}
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}
	parts := make([]string, 0, len(validationErrs))
	for _, e := range validationErrs {
		parts = append(parts, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
	}
	return fmt.Errorf("validation failed: %s", strings.Join(parts, "; "))
}
