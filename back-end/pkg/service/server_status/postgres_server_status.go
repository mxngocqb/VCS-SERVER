package service

import (
	"fmt"

	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

// ServerRepository holds DB connection logic
type ServerRepository struct {
	DB *gorm.DB
}

// NewServerRepository creates a new instance of ServerRepository.
func NewServerRepository(db *gorm.DB) *ServerRepository {
	return &ServerRepository{DB: db}
}


func New(config *config.Config) (*DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.DB.Host, config.DB.User, config.DB.Password, config.DB.Name, config.DB.Port)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	} 

	return &DB{db}, err
}

// Update updates a server by id
func (ss *ServerRepository) Update(id string, s *model.Server) error {
	updateFields := make(map[string]interface{})

	updateFields["status"] = s.Status
	
    if s.Name != "" {
        updateFields["name"] = s.Name
    }

    if s.IP != "" {
        updateFields["ip"] = s.IP
    }

    if len(updateFields) == 0 {
        return nil // No fields to update
    }

	return ss.DB.Model(&model.Server{}).Where("id = ?", id).Updates(updateFields).Error
}