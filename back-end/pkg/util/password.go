package util

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)


// Hash a password with a cost of 14 using bcrypt
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	// Generate a hashed password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// Check if the provided password matches the hashed password
func CheckPasswordHash(password, hash string) error {
	// Compare the provided password with the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		// Password does not match
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("invalid username or password")
		}
		return err
	}
	return nil
}
