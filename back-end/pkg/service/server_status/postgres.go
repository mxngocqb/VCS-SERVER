package service

import (
	"fmt"
	"log"

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

	log.Println("Connected to database")

	return &DB{db}, err
}

// Update updates the status of a server by ID
func (ss *ServerRepository) Update(id string, status bool) error {
	return ss.DB.Model(&model.Server{}).Where("id = ?", id).Update("status", status).Error
}


// GetServerByID retrieves a server by ID.
func (ss *ServerRepository) GetServerByID(id string) (*model.Server, error) {
	var server model.Server
	err := ss.DB.First(&server, id).Error
	return &server, err
}
