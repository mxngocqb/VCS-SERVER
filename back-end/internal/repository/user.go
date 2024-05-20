package repository

import (
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
)

type UserRepository	interface {
	Create(u *model.User) error
	Update(u *model.User) error
	Delete(u *model.User) error
	GetRoles(roleIDs []uint) ([]model.Role, error)
	GetUserRoleIDs(userID uint) ([]uint, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(id string) (*model.User, error)
}