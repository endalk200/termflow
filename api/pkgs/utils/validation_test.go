package utils_test

import (
	"testing"

	"github.com/endalk200/termflow-api/pkgs/utils"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var validate *validator.Validate

// TestPayload represents a struct with various validation rules
type TestPayload struct {
	Command  string `validate:"required"`               // Required field
	TagID    string `validate:"required"`               // Required field
	Email    string `validate:"required,email"`         // Email validation
	Age      int    `validate:"required,gte=18,lte=60"` // Age between 18 and 60
	Password string `validate:"required,min=8"`         // Minimum length of 8
	Website  string `validate:"required,url"`           // URL validation
}

func init() {
	validate = validator.New() // Initialize validator
}

func TestValidatePayload(t *testing.T) {
	// Test case: All fields missing
	t.Run("Missing Fields", func(t *testing.T) {
		payload := TestPayload{}
		customErrors, err := utils.ValidateAndFormatErrors(payload)

		assert.Nil(t, err)
		assert.NotNil(t, customErrors)

		// Check specific error messages, ensuring the field exists and has at least one error
		if errors, ok := customErrors.FieldErrors["Command"]; ok && len(errors) > 0 {
			assert.Equal(t, "Field Command must be required", errors[0]["required"])
		}
		if errors, ok := customErrors.FieldErrors["TagID"]; ok && len(errors) > 0 {
			assert.Equal(t, "Field TagID must be required", errors[0]["required"])
		}
		if errors, ok := customErrors.FieldErrors["Email"]; ok && len(errors) > 0 {
			assert.Equal(t, "Field Email must be required", errors[0]["required"])
		}
		if errors, ok := customErrors.FieldErrors["Age"]; ok && len(errors) > 0 {
			assert.Equal(t, "Field Age must be required", errors[0]["required"])
		}
		if errors, ok := customErrors.FieldErrors["Password"]; ok && len(errors) > 0 {
			assert.Equal(t, "Field Password must be required", errors[0]["required"])
		}
		if errors, ok := customErrors.FieldErrors["Website"]; ok && len(errors) > 0 {
			assert.Equal(t, "Field Website must be required", errors[0]["required"])
		}
	})

	// Test case: Invalid email and age out of range
	t.Run("Invalid Email and Age", func(t *testing.T) {
		payload := TestPayload{
			Command:  "some-command",
			TagID:    "some-tag-id",
			Email:    "invalid-email", // Invalid email format
			Age:      17,              // Below minimum age
			Password: "validPassword",
			Website:  "http://validurl.com",
		}
		customErrors, err := utils.ValidateAndFormatErrors(payload)

		assert.Nil(t, err)
		assert.NotNil(t, customErrors)

		// Check specific error messages
		if errors, ok := customErrors.FieldErrors["Email"]; ok && len(errors) > 0 {
			assert.Equal(t, "Field Email must be email", errors[0]["email"])
		}
		if errors, ok := customErrors.FieldErrors["Age"]; ok && len(errors) > 0 {
			assert.Equal(t, "Field Age must be gte", errors[0]["gte"])
		}
	})

	// Test case: Password too short and invalid URL
	t.Run("Short Password and Invalid URL", func(t *testing.T) {
		payload := TestPayload{
			Command:  "some-command",
			TagID:    "some-tag-id",
			Email:    "valid@example.com",
			Age:      30,
			Password: "short",       // Too short
			Website:  "invalid-url", // Invalid URL
		}
		customErrors, err := utils.ValidateAndFormatErrors(payload)

		assert.Nil(t, err)
		assert.NotNil(t, customErrors)

		// Check specific error messages
		if errors, ok := customErrors.FieldErrors["Password"]; ok && len(errors) > 0 {
			assert.Equal(t, "Field Password must be min", errors[0]["min"])
		}
		if errors, ok := customErrors.FieldErrors["Website"]; ok && len(errors) > 0 {
			assert.Equal(t, "Field Website must be url", errors[0]["url"])
		}
	})

	// Test case: Valid payload (all fields present and valid)
	t.Run("Valid Payload", func(t *testing.T) {
		payload := TestPayload{
			Command:  "some-command",
			TagID:    "some-tag-id",
			Email:    "valid@example.com",
			Age:      30,
			Password: "validPassword",
			Website:  "http://validurl.com",
		}
		customErrors, err := utils.ValidateAndFormatErrors(payload)

		assert.Nil(t, err)          // Validation should pass
		assert.Nil(t, customErrors) // No errors expected
	})
}
