package auth_test

import (
	"testing"

	"github.com/endalk200/termflow-api/pkgs/auth"
)

func TestHashPassword(t *testing.T) {
	password := "my_secure_password"

	// Test bcrypt hashing
	bcryptHash, err := auth.HashPassword(password, auth.Bcrypt)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if bcryptHash == "" {
		t.Fatal("Expected bcrypt hash to be non-empty")
	}

	// Test SHA256 hashing
	sha256Hash, err := auth.HashPassword(password, auth.SHA256)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if sha256Hash == "" {
		t.Fatal("Expected SHA256 hash to be non-empty")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "my_secure_password"

	// Hash the password using bcrypt
	bcryptHash, err := auth.HashPassword(password, auth.Bcrypt)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check bcrypt hash
	isValid, err := auth.CompareHash(password, bcryptHash, auth.Bcrypt)
	if err != nil || !isValid {
		t.Fatalf("Expected valid password check for bcrypt, got error %v", err)
	}

	// Check with an incorrect password
	isValid, err = auth.CompareHash("wrong_password", bcryptHash, auth.Bcrypt)
	if isValid {
		t.Fatal("Expected invalid password check for bcrypt")
	}

	// Hash the password using SHA256
	sha256Hash, err := auth.HashPassword(password, auth.SHA256)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check SHA256 hash
	isValid, err = auth.CompareHash(password, sha256Hash, auth.SHA256)
	if err != nil || !isValid {
		t.Fatalf("Expected valid password check for SHA256, got error %v", err)
	}

	// Check with an incorrect password
	isValid, err = auth.CompareHash("wrong_password", sha256Hash, auth.SHA256)
	if isValid {
		t.Fatal("Expected invalid password check for SHA256")
	}
}

func TestUnsupportedHashAlgorithm(t *testing.T) {
	_, err := auth.HashPassword("password", auth.HashAlgorithm(99)) // Using an unsupported algorithm
	if err == nil {
		t.Fatal("Expected an error for unsupported hash algorithm, got none")
	}

	_, err = auth.CompareHash("password", "hash", auth.HashAlgorithm(99)) // Using an unsupported algorithm
	if err == nil {
		t.Fatal("Expected an error for unsupported hash algorithm, got none")
	}
}
