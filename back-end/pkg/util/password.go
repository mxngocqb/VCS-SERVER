package util

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a plain text password and returns a bcrypt hashed version.
// It returns an error if the password string is empty.
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// CheckPasswordHash compares a bcrypt hashed password with its possible plaintext equivalent
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		// Password does not match
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("invalid username or password")
		}
		// Other error
		return err
	}

	// Credentials are valid
	return nil
}
