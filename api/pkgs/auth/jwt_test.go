package auth_test

import (
	"os"
	"testing"
	"time"

	"github.com/endalk200/termflow-api/pkgs/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const testPrivateKeyPath = "../../test/test_private_key.pem"
const testPublicKeyPath = "../../test/test_public_key.pem"

// Test LoadPrivateKey
func TestLoadPrivateKey(t *testing.T) {
	privateKey, err := auth.LoadPrivateKey(testPrivateKeyPath)
	assert.NoError(t, err, "Loading private key should not return an error")
	assert.NotNil(t, privateKey, "Private key should not be nil")
}

// Test LoadPublicKey
func TestLoadPublicKey(t *testing.T) {
	publicKey, err := auth.LoadPublicKey(testPublicKeyPath)
	assert.NoError(t, err, "Loading public key should not return an error")
	assert.NotNil(t, publicKey, "Public key should not be nil")
}

// Test GenerateJWT
func TestGenerateJWT(t *testing.T) {
	// Generate test claims
	claims := jwt.RegisteredClaims{
		Issuer:    "Company",
		Subject:   "1234567890",
		Audience:  jwt.ClaimStrings{"test-audience"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	tokenString, err := auth.GenerateJWT(claims)
	assert.NoError(t, err, "Generating JWT should not return an error")
	assert.NotEmpty(t, tokenString, "Generated token should not be empty")
}

// Test VerifyJWT - valid token
func TestVerifyJWT(t *testing.T) {
	// Create valid JWT
	claims := jwt.RegisteredClaims{
		Issuer:    "Company",
		Subject:   "1234567890",
		Audience:  jwt.ClaimStrings{"test-audience"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	tokenString, err := auth.GenerateJWT(claims)
	assert.NoError(t, err, "Generating valid JWT should not return an error")

	// Verify JWT
	token, err := auth.VerifyJWT(tokenString)
	assert.NoError(t, err, "Verifying valid JWT should not return an error")
	assert.True(t, token.Valid, "Token should be valid")

	// Verify the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		assert.Equal(t, "Company", claims["iss"], "Issuer should be 'Company'")
	} else {
		t.Fatalf("Token claims are not in expected format")
	}
}

// Test VerifyJWT - invalid token
func TestVerifyJWT_Invalid(t *testing.T) {
	invalidTokenString := "invalid.token.string"

	_, err := auth.VerifyJWT(invalidTokenString)
	assert.Error(t, err, "Verifying an invalid token should return an error")
}

// Test LoadPrivateKey - invalid file path
func TestLoadPrivateKey_InvalidPath(t *testing.T) {
	_, err := auth.LoadPrivateKey("invalid_private_key_path.pem")
	assert.Error(t, err, "Loading private key from invalid path should return an error")
}

// Test LoadPublicKey - invalid file path
func TestLoadPublicKey_InvalidPath(t *testing.T) {
	_, err := auth.LoadPublicKey("invalid_public_key_path.pem")
	assert.Error(t, err, "Loading public key from invalid path should return an error")
}

// Cleanup function
func cleanupTestFiles() {
	_ = os.Remove(testPrivateKeyPath)
	_ = os.Remove(testPublicKeyPath)
}
