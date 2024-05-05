package internal

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/handler"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/repository"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/util"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Start initializes and starts the Echo API server
func Start(cfg *config.Config) error {
	// Initialize the database
	db, err := repository.New(cfg)
	if err != nil {
		return err
	} else {
		log.Printf("Connected to Postgres")
	}
	db.Config.Logger = db.Config.Logger.LogMode(4)
	// Initialize Redis service
	util.InitRedis()
	// Initialize Elastic service
	elasticService := util.NewElasticsearch()
	if err := elasticService.CreateStatusLogIndex(); err != nil {
		return err
	}

	// Initialize stores
	userRepository := repository.NewUserRepository(db.DB)
	// serverStore := repository.NewServerRepository(db.DB)

	// Initialize services
	handler.NewRbacService(userRepository)

	// Set up Echo Server
	e := echo.New()
	//e.HideBanner = true
	//e.HidePort = true
	//
	//// Configure lumberjack logger
	//e.Logger.SetOutput(util.LogConfig)
	//
	//// Middleware to log HTTP requests
	//e.Use(middleware.Logger(), middleware.Recover())

	// Set up Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// New Create user endpoint

	// Closed group

	// Start the server

	// Schedule daily report

	return nil
}
