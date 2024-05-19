package transport

import (
	"net/http"

	_ "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/handler/user"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
)

// HTTP represents the HTTP transport for user endpoints.
type HTTP struct {
	service user.IUserService
}

// NewHTTP creates a new HTTP Group handler for user endpoints.
func NewHTTP(r *echo.Group, service *user.Service) {
	h := HTTP{service}

	ur := r.Group("/users")
	ur.GET("/:id", h.View)
	ur.POST("", h.Create)
	ur.PUT("/:id", h.Update)
	ur.DELETE("/:id", h.Delete)
}

// View godoc
// @Summary View user
// @Description Get details of a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} echo.HTTPError "Bad Request - Invalid user ID format"
// @Failure 404 {object} echo.HTTPError "Not Found - User not found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error - Unable to retrieve user"
// @Security Bearer
// @Router /users/{id} [get]
func (h HTTP) View(c echo.Context) error {
	id := c.Param("id")
	viewUser, err := h.service.View(c, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, viewUser)
}

// Create godoc
// @Summary Create user
// @Description Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body CreateRequest true "Create User"
// @Success 201 {object} model.User
// @Failure 400 {object} echo.HTTPError "Bad Request - Invalid user data"
// @Failure 403 {object} echo.HTTPError "Forbidden - Insufficient permissions"
// @Failure 500 {object} echo.HTTPError "Internal Server Error - Unable to create user"
// @Security Bearer
// @Router /users [post]
func (h HTTP) Create(c echo.Context) error {
	r := new(CreateRequest)
	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// Create a new user model
	newUser := &model.User{
		Username: r.Username,
		Password: r.Password,
		RoleIDs:  r.RoleIDs,
	}

	// Create new user with roles in the database
	createdUser, err := h.service.Create(c, newUser)
	if err != nil {
		return err 
	}

	return c.JSON(http.StatusCreated, createdUser)
}

// Update godoc
// @Summary Update user
// @Description Update a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body UpdateRequest true "Update User"
// @Success 200 {object} model.User
// @Failure 400 {object} echo.HTTPError "Bad Request - Invalid user data or ID"
// @Failure 403 {object} echo.HTTPError "Forbidden - Insufficient permissions"
// @Failure 404 {object} echo.HTTPError "Not Found - User not found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error - Unable to update user"
// @Security Bearer
// @Router /users/{id} [put]
func (h HTTP) Update(c echo.Context) error {
	r := new(UpdateRequest)
	id := c.Param("id")
	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Create a new user model
	updatedUser := &model.User{
		Username: r.Username,
		Password: r.Password,
		RoleIDs:  r.RoleIDs,
	}

	// Update user with roles in the database
	updatedUser, err := h.service.Update(c, id, updatedUser)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, updatedUser)
}

// Delete godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 403 {object} echo.HTTPError "Forbidden - Insufficient permissions"
// @Failure 404 {object} echo.HTTPError "Not Found - User not found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error - Unable to delete user"
// @Router /users/{id} [delete]
// @Security Bearer
func (h HTTP) Delete(c echo.Context) error {
	id := c.Param("id")
	err := h.service.Delete(c, id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}


