package model

// Operation defines an operation that can be performed
type Operation struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"uniqueIndex" json:"name"`
}
