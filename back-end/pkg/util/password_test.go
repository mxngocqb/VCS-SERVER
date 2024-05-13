package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "admin"

	// Test hashing the password
	hashedPassword, err := HashPassword(password)
	t.Logf(hashedPassword)
	assert.NoError(t, err, "Hashing the password should not produce an error")
	assert.NotEqual(t, password, hashedPassword, "Hashed password should be different from the plaintext password")

	// Test checking the password hash
	err = CheckPasswordHash(password, hashedPassword)
	assert.NoError(t, err, "Checking the correct password against its hash should not produce an error")

	// Test checking the password hash with a wrong password
	wrongPassword := "wrongpassword"
	err = CheckPasswordHash(wrongPassword, hashedPassword)
	assert.Error(t, err, "Checking the wrong password against its hash should produce an error")
	assert.Equal(t, "invalid username or password", err.Error(), "Error message should be specific for mismatched password")
}

func TestHashPasswordFailure(t *testing.T) {
	// This test assumes that there might be an error in the bcrypt hashing process
	// which is unlikely to be triggered unless there is a specific error like an empty password

	emptyPassword := ""

	// Test hashing an empty password
	_, err := HashPassword(emptyPassword)
	assert.Error(t, err, "Hashing an empty password should produce an error")
}
