package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type HashAlgorithm int

const (
	Bcrypt HashAlgorithm = iota
	SHA256
)

func HashPassword(raw string, algorithm HashAlgorithm) (string, error) {
	switch algorithm {
	case Bcrypt:
		bytes, err := bcrypt.GenerateFromPassword([]byte(raw), 14)
		return string(bytes), err
	case SHA256:
		hasher := sha256.New()
		hasher.Write([]byte(raw))
		return hex.EncodeToString(hasher.Sum(nil)), nil
	default:
		return "", fmt.Errorf("unsupported hash algorithm")
	}
}

func CompareHash(raw, hash string, algorithm HashAlgorithm) (bool, error) {
	switch algorithm {
	case Bcrypt:
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw))
		return err == nil, err
	case SHA256:
		hashedInput := sha256.Sum256([]byte(raw))
		return hex.EncodeToString(hashedInput[:]) == hash, nil
	default:
		return false, fmt.Errorf("unsupported hash algorithm")
	}
}
