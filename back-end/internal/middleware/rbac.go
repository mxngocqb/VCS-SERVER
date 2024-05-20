package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	customjwt "github.com/mxngocqb/VCS-SERVER/back-end/pkg/util"
) 

// RoleMiddleware sets up the role-based access control middleware.
func RoleMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*customjwt.UserClaims)

			c.Set("id", claims.ID)
			c.Set("username", claims.Username)

			return next(c)
		}
	}
}
