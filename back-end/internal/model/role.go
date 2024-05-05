package model

// Role defines a role in the system
type Role struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"uniqueIndex" json:"name"`
}

// RolePermissions defines the permissions a role has on an operation for a resource
type RolePermissions struct {
	RoleID      uint `gorm:"primaryKey" json:"role_id"`
	OperationID uint `gorm:"primaryKey" json:"operation_id"`
	ResourceID  uint `gorm:"primaryKey" json:"resource_id"`
}
