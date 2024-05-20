package transport

// CreateRequest defines the request payload for creating a new user.
type CreateRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	RoleIDs  []uint `json:"role_ids" validate:"required,min=1"`
}

// UpdateRequest defines the request payload for updating a user.
type UpdateRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	RoleIDs  []uint `json:"role_ids" validate:"required"`
}
