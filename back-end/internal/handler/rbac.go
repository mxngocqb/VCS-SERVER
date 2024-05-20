package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/repository"
)

type RbacService interface {
	EnforceRole(c echo.Context, roleID uint) error
}

type RbacServiceImpl struct {
	repository repository.UserRepository
}

func NewRbacServiceImpl(repository repository.UserRepository) *RbacServiceImpl {
	return &RbacServiceImpl{
		repository: repository,
	}
}

func (r *RbacServiceImpl) EnforceRole(c echo.Context, roleID uint) error {
	// Get the user's role from the context
	userID := c.Get("id").(uint)

	roleIDs, err := r.repository.GetUserRoleIDs(userID)
	if err != nil {
		return err
	}

	// Check if the user has the required role
	for _, id := range roleIDs {
		if id == roleID {
			return nil
		}
	}

	return echo.NewHTTPError(403, "Forbidden - Insufficient permissions")
}
