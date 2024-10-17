package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CustomValidationErrors struct {
	FieldErrors map[string][]map[string]string `json:"field_name"`
}

// validateAndFormatErrors validates the struct and formats the errors.
func ValidateAndFormatErrors(payload interface{}) (*CustomValidationErrors, error) {
	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		validationErrors := &CustomValidationErrors{
			FieldErrors: make(map[string][]map[string]string),
		}

		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.StructField()
			errorType := err.Tag()
			message := fmt.Sprintf("Field %s must be %s", fieldName, err.Tag())

			// Append the error to the field
			validationErrors.FieldErrors[fieldName] = append(validationErrors.FieldErrors[fieldName], map[string]string{
				errorType: message,
			})
		}
		return validationErrors, nil
	}
	return nil, nil
}
