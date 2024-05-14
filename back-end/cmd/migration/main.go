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
	if err != nil {
		panic("Failed to drop tables")
	}

	// Automatically migrate your schema.
	err = db.AutoMigrate(&model.User{}, &model.Role{}, &model.Operation{}, &model.Resource{}, &model.RolePermissions{}, &model.Server{})
	if err != nil {
		panic("Failed to migrate database")
	}
	fmt.Println("Database migrated successfully!")

	// Seed the database with roles
	if err := seedRoles(db); err != nil {
		log.Fatalf("Failed to seed roles: %v", err)
	}

	// SEed the database with users
	if err := seedUsers(db); err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	fmt.Println("Database seeded successfully!")
}

// seedRoles adds predefined roles to the roles table
func seedRoles(db *gorm.DB) error {
	// SQL statements for seeding roles
	queries := []string{
		"INSERT INTO roles (id, name) VALUES (DEFAULT, 'Admin');",
		"INSERT INTO roles (id, name) VALUES (DEFAULT, 'User');",
	}

	// Execute each query
	for _, query := range queries {
		result := db.Exec(query)
		if result.Error != nil {
			return fmt.Errorf("failed to insert role: %v", result.Error)
		}
		fmt.Printf("Inserted role with result: %v\n", result.RowsAffected)
	}

	return nil
}

// seedUsers adds predefined users to the users table
func seedUsers(db *gorm.DB) error {
    // SQL statements for seeding users
    queries := []string{
        "INSERT INTO users (username, password) VALUES ('admin', '$2a$14$NcWnkwHlM0Vmu6DeWn0afu.Pc1HbPgghQpqOB/43Ah4IJeR/UXrEW');",
        "INSERT INTO user_roles VALUES ('1','1');",
    }

    // Execute each query
    for _, query := range queries {
        result := db.Exec(query)
        if result.Error != nil {
            return fmt.Errorf("failed to insert user: %v", result.Error)
        }
        fmt.Printf("Inserted user with result: %v\n", result.RowsAffected)
    }

    return nil
}
