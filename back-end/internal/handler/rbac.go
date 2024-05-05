package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/repository"
)

type RbacService struct {
	repository *repository.UserRepository
}

func NewRbacService(repository *repository.UserRepository) *RbacService {
	return &RbacService{
		repository: repository,
	}
}

func (r *RbacService) EnforceRole(c echo.Context, roleID uint) error {

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

	return echo.NewHTTPError(403, "forbidden")
}
