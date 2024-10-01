package auth_test

import (
	"testing"

	"github.com/endalk200/termflow-api/pkgs/auth"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"

	// Test that hashing a password does not return an error
	hashedPassword, err := auth.HashPassword(password)
	assert.NoError(t, err, "Error hashing password should be nil")

	// Ensure that the hash is not the same as the plain password
	assert.NotEqual(t, hashedPassword, password, "Hashed password should not match the original password")

	// Test that the hashed password can be verified against the original password
	isValid, err := auth.CheckPasswordHash(password, hashedPassword)
	assert.NoError(t, err, "Error checking password hash should be nil")
	assert.True(t, isValid, "Password and hash should match")
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mypassword"

	// Manually hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err, "Error generating bcrypt hash should be nil")

	// Test that a correct password returns true
	isValid, err := auth.CheckPasswordHash(password, string(hashedPassword))
	assert.NoError(t, err, "Error checking password hash should be nil")
	assert.True(t, isValid, "Expected password to be valid")

	// Test that an incorrect password returns false
	incorrectPassword := "wrongpassword"
	isValid, err = auth.CheckPasswordHash(incorrectPassword, string(hashedPassword))
	assert.NoError(t, err, "Error checking password hash should be nil")
	assert.False(t, isValid, "Expected password to be invalid")

	// Additional tests for edge cases
	t.Run("EmptyPassword", func(t *testing.T) {
		// Hashing an empty password
		emptyHash, err := auth.HashPassword("")
		assert.NoError(t, err, "Error hashing empty password should be nil")

		// Verifying against the empty password
		isValid, err = auth.CheckPasswordHash("", emptyHash)
		assert.NoError(t, err, "Error checking empty password hash should be nil")
		assert.True(t, isValid, "Expected empty password to match its hash")
	})

	t.Run("EmptyHash", func(t *testing.T) {
		// Checking a password against an empty hash
		isValid, err := auth.CheckPasswordHash(password, "")
		assert.Error(t, err, "Expected an error when checking against an empty hash")
		assert.False(t, isValid, "Expected password to be invalid against an empty hash")
	})
}
