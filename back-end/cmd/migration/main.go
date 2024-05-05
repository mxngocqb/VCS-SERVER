package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Collect database connection information
	dbHost, dbUser, dbPassword, dbName, dbPort :=
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT")

	// Check if all environment variables are present
	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
		log.Fatal("Database environment variables are not all set")
	}

	// Construct the DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Drop all tables and recreate them.
	err = db.Migrator().DropTable(&model.User{}, &model.Role{}, &model.Operation{}, &model.Resource{}, &model.RolePermissions{}, &model.Server{})

	// Automatically migrate your schema.
	err = db.AutoMigrate(&model.User{}, &model.Role{}, &model.Operation{}, &model.Resource{}, &model.RolePermissions{}, &model.Server{})
	if err != nil {
		panic("Failed to migrate database")
	}

	fmt.Println("Database migrated successfully!")
}
