package user

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/handler"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/repository"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/util"
)

type IUserService interface {
	View(ctx echo.Context, id string) (*model.User, error)
	Create(ctx echo.Context, u *model.User) (*model.User, error)
	Update(ctx echo.Context, id string, u *model.User) (*model.User, error)
	Delete(ctx echo.Context, id string) error
}
type Service struct {
	repository repository.UserRepository
	rbac  handler.RbacService
}

// NewUserService creates a new instance of UserService.
func NewUserService(repository repository.UserRepository, rbac handler.RbacService) *Service {
	return &Service{
		repository: repository,
		rbac:  rbac,
	}
}

// View retrieves a user by ID.
func (s *Service) View(ctx echo.Context, id string) (*model.User, error) {
	return s.repository.GetUserByID(id)
}

// Create creates a new user in the database.
func (s *Service) Create(c echo.Context, u *model.User) (*model.User, error) {
	requiredRoleID := uint(1)
	if err := s.rbac.EnforceRole(c, requiredRoleID); err != nil {
		return nil, err
	}

	// Hash the password before saving it
	hashedPassword, err := util.HashPassword(u.Password)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	u.Password = hashedPassword

	// Fetch roles based on provided role IDs and assign them to the user
	roles, err := s.repository.GetRoles(u.RoleIDs)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	u.Roles = roles

	// Create the user in the database
	if err := s.repository.Create(u); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return u, nil
}

// Update updates a user in the database using the user id.
func (s *Service) Update(c echo.Context, id string, u *model.User) (*model.User, error) {
	requiredRoleID := uint(1)
	if err := s.rbac.EnforceRole(c, requiredRoleID); err != nil {
		return nil, err
	}

	// Retrieve the existing user by ID
	existingUser, err := s.repository.GetUserByID(id)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// If no user is found with the given ID, return an error
	if existingUser == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	// Hash the new password before saving it
	if u.Password != "" && u.Password != existingUser.Password {
		hashedPassword, err := util.HashPassword(u.Password)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		existingUser.Password = hashedPassword
	}

	// Update roles if provided
	if len(u.RoleIDs) > 0 {
		roles, err := s.repository.GetRoles(u.RoleIDs)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		existingUser.Roles = roles
	}

	existingUser.Username = u.Username

	// Update the user in the database
	if err := s.repository.Update(existingUser); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return existingUser, nil
}

// Delete deletes a user from the database.
func (s *Service) Delete(c echo.Context, id string) error {
	requiredRoleID := uint(1)
	if err := s.rbac.EnforceRole(c, requiredRoleID); err != nil {
		return err
	}

	// Retrieve the existing user by ID
	existingUser, err := s.repository.GetUserByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// If no user is found with the given ID, return an error
	if existingUser == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	// Delete the user in the database
	if err := s.repository.Delete(existingUser); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
