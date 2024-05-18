package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/repository"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/util"
)

type IService interface {
	Authenticate(username, password string) (*model.User, error)
}

// Service provides user authentication services.
type Service struct {
	repository *repository.UserRepository
}

// NewAuthService creates a new authentication service.
// @Summary Create authentication service
// @Description It prepares a new instance of the AuthService which interacts with the user repository.
// @Tags Authentication
func NewAuthService(repository *repository.UserRepository) *Service {
	return &Service{
		repository: repository,
	}
}

// Authenticate checks if user credentials are valid.
// @Summary Authenticate user
// @Description Checks if the provided username and password are correct.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Param password path string true "Password"
// @Success 200 {object} model.User "Authenticated user"
// @Failure 400 {object} string "User not found"
// @Failure 401 {object} string "Invalid username or password"
func (s *Service) Authenticate(username, password string) (*model.User, error) {
	// Retrieve the user by username
	user, err := s.repository.GetUserByUsername(username)
	if err != nil {
		return &model.User{}, err
	}

	// Compare the provided password with the hashed password
	err = util.CheckPasswordHash(password, user.Password)
	if err != nil {
		return &model.User{}, echo.NewHTTPError(401, "invalid username or password")
	}

	// Credentials are valid
	return user, nil
}
