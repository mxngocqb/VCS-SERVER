package repository

import (
	"fmt"
	"log"
	"sync"

	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

var (
    db   *gorm.DB
    once sync.Once
)

// Using singleton pattern to connect to the database
func InitDB(config *config.Config, logger logger.Interface)  {
	once.Do(func() {
		var err error
		// Create a new connection to the database
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			config.DB.Host, config.DB.User, config.DB.Password, config.DB.Name, config.DB.Port)

		// Connect to the database
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger: logger,
		})
		if err != nil {
			log.Printf("Error connecting to database: %v", err)
		} 

	})
}

func GetDB() (*DB, error) {
    if db == nil {
        return nil, fmt.Errorf("db is nil")
    }
    return &DB{db}, nil
}