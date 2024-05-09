package model

import (
	"gorm.io/gorm"
)

// Server defines a server in the system
type Server struct {
	gorm.Model
	Name   string `json:"name" gorm:"type:varchar(255);not null"`
	Status bool   `json:"status" gorm:"not null"`
	IP     string `json:"ip" gorm:"type:varchar(15);not null;unique"`
}


type ServerSwag struct {
	Name   string `json:"name" gorm:"type:varchar(255);not null"`
	Status bool   `json:"status" gorm:"not null"`
	IP     string `json:"ip" gorm:"type:varchar(15);not null;unique"`
}
