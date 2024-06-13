package repository

import (
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"gorm.io/gorm"
)


// UserRepositoryImpl holds DB connection logic
type UserRepositoryImpl struct {
	DB *gorm.DB
}

// NewUserRepositoryImpl creates a new instance of UserRepositoryImpl.
func NewUserRepositoryImpl(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}


// Create creates a new user in the database.
func (us *UserRepositoryImpl) Create(u *model.User) error {
	return us.DB.Create(u).Error
}

// GetUsers retrieves all users from the database.
func (us *UserRepositoryImpl) GetUsers() ([]model.User, error) {
	var users []model.User
	err := us.DB.Find(&users).Error
	return users, err
}

// Update updates a user in the database.
func (us *UserRepositoryImpl) Update(u *model.User) error {
	return us.DB.Save(u).Error
}

// Delete deletes a user from the database.
func (us *UserRepositoryImpl) Delete(u *model.User) error {
	return us.DB.Delete(u).Error
}

// GetRoles retrieves roles from the database.
func (us *UserRepositoryImpl) GetRoles(roleIDs []uint) ([]model.Role, error) {
	var roles []model.Role
	err := us.DB.Find(&roles, roleIDs).Error
	return roles, err
}

// GetUserRoleIDs fetches all role IDs associated with a given user ID.
func (us *UserRepositoryImpl) GetUserRoleIDs(userID uint) ([]uint, error) {
	var user model.User
	var roleIDs []uint

	// Fetch the user along with their roles
	result := us.DB.Preload("Roles").First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}

	// Extract the role IDs from the user's roles
	for _, role := range user.Roles {
		roleIDs = append(roleIDs, role.ID)
	}

	return roleIDs, nil
}

// GetUserByUsername finds a user by their username and preloads the roles.
func (us *UserRepositoryImpl) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	result := us.DB.Preload("Roles").Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByID finds a user by their ID and preloads the roles.
func (us *UserRepositoryImpl) GetUserByID(id string) (*model.User, error) {
	var user model.User
	result := us.DB.Preload("Roles").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
