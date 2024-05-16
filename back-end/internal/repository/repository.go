package repository

import (
	"fmt"

	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

func New(config *config.Config, logger logger.Interface) (*DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.DB.Host, config.DB.User, config.DB.Password, config.DB.Name, config.DB.Port)

	
	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: logger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	} 

	return &DB{db}, err
}
