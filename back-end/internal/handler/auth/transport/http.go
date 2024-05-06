package transport

import (
	"github.com/labstack/echo/v4"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/handler/auth"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/util"
	"net/http"
)

// HTTP represents the HTTP transport for the auth service.
type HTTP struct {
	service auth.IService
}

// NewHTTP sets up routes related to authentication services
// @Summary Set up authentication routes
// @Description It initializes routes for login and other authentication-related processes.
// @Tags Authentication
// @Accept json
// @Produce json
func NewHTTP(r *echo.Group, service *auth.Service) {
	h := HTTP{service}

	r.POST("/login", h.Login)
}

// Login authenticates a user via username and password.
// @Summary User login
// @Description Authenticates user and returns a JWT token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param LoginRequest body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} string "Invalid input"
// @Failure 500 {object} string "Internal Server Error"
// @Router /login [post]
func (h HTTP) Login(c echo.Context) error {
	r := new(LoginRequest)
	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: "+err.Error())
	}
	// Authenticate the user
	user, err := h.service.Authenticate(r.Username, r.Password)
	if err != nil {
		return err
	}

	token, err := util.GenerateToken(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}
