package model

// Resource defines a resource upon which operations can be performed
type Resource struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"uniqueIndex" json:"name"`
}
