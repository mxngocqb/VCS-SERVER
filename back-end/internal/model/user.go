package model

import (
	"gorm.io/gorm"
)

// User defines a user in the system with associated roles
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex" json:"username"`
	Password string `json:"password"` // Password should be hashed and never returned in API calls
	Roles    []Role `gorm:"many2many:user_roles;" json:"roles"`
	RoleIDs  []uint `gorm:"-" json:"role_ids"`
}

