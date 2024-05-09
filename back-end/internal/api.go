package internal

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/handler"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/handler/user"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/handler/auth"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/handler/server"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/repository"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/util"
	"github.com/go-playground/validator/v10"
	at "github.com/mxngocqb/VCS-SERVER/back-end/internal/handler/auth/transport"
	st "github.com/mxngocqb/VCS-SERVER/back-end/internal/handler/server/transport"
	ut "github.com/mxngocqb/VCS-SERVER/back-end/internal/handler/user/transport"
	echoSwagger "github.com/swaggo/echo-swagger"
	echojwt "github.com/labstack/echo-jwt/v4"
	custommiddleware "github.com/mxngocqb/VCS-SERVER/back-end/internal/middleware"
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

	// Initialize Repos 
	userRepository := repository.NewUserRepository(db.DB)
	serverRepository := repository.NewServerRepository(db.DB)

	// Initialize services
	rbacService := handler.NewRbacService(userRepository)
	userService := user.NewUserService(userRepository, rbacService)
	authService := auth.NewAuthService(userRepository)
	serverService := server.NewServerService(serverRepository, rbacService, elasticService)
	
	user.NewUserService(userRepository, rbacService)
	auth.NewAuthService(userRepository)
	user.NewUserService(userRepository, rbacService)
	// Set up Echo Server
	e := echo.New()
	//e.HideBanner = true
	//e.HidePort = true
	//
	//// Configure lumberjack logger
	//e.Logger.SetOutput(util.LogConfig)
	//s
	//// Middleware to log HTTP requests
	//e.Use(middleware.Logger(), middleware.Recover())
	e.Validator = &util.CustomValidator{Validator: validator.New()}
	e.Binder = &util.CustomBinder{Binder: &echo.DefaultBinder{}}
	// Set up Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	v1 := e.Group("/api/v1")
	// New Create user endpoint
	at.NewHTTP(v1, authService)
	// jwtBlocked group
	jwtBlocked := v1.Group("")
	jwtBlocked.Use(echojwt.WithConfig(custommiddleware.JWTMiddleware()))
	jwtBlocked.Use(custommiddleware.RoleMiddleware())
	ut.NewHTTP(jwtBlocked, userService)
	st.NewHTTP(jwtBlocked, serverService)
	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
	// Schedule daily report
	util.ScheduleDailyReport()
	return nil
}
