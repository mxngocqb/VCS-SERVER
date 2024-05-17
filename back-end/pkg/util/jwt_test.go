package util

import (
	"testing"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGenerateToken(t *testing.T) {
	// Mocking time for consistent test results
	now := time.Now()

	user := &model.User{
		Model:    gorm.Model{ID: 1},
		Username: "testuser",
	}

	tokenString, err := GenerateToken(user)
	assert.Nil(t, err, "Token generation should not return an error")
	assert.NotEmpty(t, tokenString, "Token should not be empty")

	// Parse the token to check if the claims are correct
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil
	})

	assert.Nil(t, err, "Token parsing should not return an error")
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		assert.Equal(t, user.ID, claims.ID, "The ID in the token should match the user's ID")
		assert.Equal(t, user.Username, claims.Username, "The username in the token should match the user's username")
		assert.WithinDuration(t, now.Add(time.Hour*72), claims.ExpiresAt.Time, time.Second, "Expiration time should be exactly 72 hours from now")
	} else {
		t.Errorf("Failed to parse claims or token is not valid")
	}
}