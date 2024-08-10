package helpers

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	issuer = "twoMatchesCorp"
)

// Load private key from private key file locally
func LoadPrivateKey(path string) (ed25519.PrivateKey, error) {
	privateKeyPEM, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %v", err)
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %v", err)
	}

	ed25519PrivateKey, ok := privateKey.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an ed25519 private key")
	}

	return ed25519PrivateKey, nil
}

func LoadPublicKey(path string) (ed25519.PublicKey, error) {
	publicKeyPEM, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading public key file: %v", err)
	}

	block, _ := pem.Decode(publicKeyPEM)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %v", err)
	}

	ed25519PublicKey, ok := publicKey.(ed25519.PublicKey)
	if !ok {

		return nil, fmt.Errorf("not an ed25519 public key")
	}

	return ed25519PublicKey, nil
}

func GenerateJWT(claims jwt.RegisteredClaims) (string, error) {
	// Load keys from files
	privateKey, err := LoadPrivateKey("private_key.pem")
	if err != nil {
		return "", fmt.Errorf("Error loading private key: %s", err)
	}

	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	publicKey, err := LoadPublicKey("public_key.pem")
	if err != nil {
		return nil, fmt.Errorf("Error loading public key: %s", err)
	}

	// Parse and verify the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate the token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["iss"] != issuer {
			return nil, fmt.Errorf("invalid issuer")
		}
	}

	return token, nil
}
