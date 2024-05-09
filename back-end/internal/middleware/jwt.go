package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	customjwt "github.com/mxngocqb/VCS-SERVER/back-end/pkg/util"
)

// JWTMiddleware creates a new JWT middleware with the specified signing key.
func JWTMiddleware() echojwt.Config {
	//Config middleware with the custom claims type
	config := echojwt.Config{
		SigningKey: []byte("your_secret_key"),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(customjwt.CustomClaims) // Use MapClaims for JWT data
		},
		SigningMethod: "HS256", // Use HS256 signing method
	}
	return config
}
