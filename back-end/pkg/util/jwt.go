package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"time"
)

type CustomClaims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token for authenticated users
func GenerateToken(u *model.User) (string, error) {
	// Convert roles to a slice of role names for inclusion in the token.

	claims := &CustomClaims{
		ID:       u.ID,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // Set the expiration time
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("your_secret_key")) // Use the same key as in the middleware
	if err != nil {
		fmt.Println("Error generating token: ", err)
		return "", err
	}

	return t, nil
}
