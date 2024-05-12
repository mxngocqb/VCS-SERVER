package internal

import (
	"log"

	"github.com/go-playground/validator/v10"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/handler"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/handler/auth"
	at "github.com/mxngocqb/VCS-SERVER/back-end/internal/handler/auth/transport"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/handler/server"
	st "github.com/mxngocqb/VCS-SERVER/back-end/internal/handler/server/transport"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/handler/user"
	ut "github.com/mxngocqb/VCS-SERVER/back-end/internal/handler/user/transport"
	custommiddleware "github.com/mxngocqb/VCS-SERVER/back-end/internal/middleware"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/repository"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/cache"
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
	redisConfig := config.NewRedisConfig(cfg)
	redisClient, expiration, err := config.ConnectRedis(redisConfig)
	serverCache := cache.NewServerRedisCache(redisClient, expiration)

	// Initialize Elastic service
	elasticService := service.NewElasticsearch()
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
	serverService := server.NewServerService(serverRepository, rbacService, elasticService, serverCache)
	// serverService := server.NewServerService(serverRepository, rbacService, elasticService)

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
	api := e.Group("/api")

	// New Create user endpoint
	at.NewHTTP(api, authService)

	// jwtBlocked group
	jwtBlocked := api.Group("")
	jwtBlocked.Use(echojwt.WithConfig(custommiddleware.JWTMiddleware()))
	jwtBlocked.Use(custommiddleware.RoleMiddleware())

	ut.NewHTTP(jwtBlocked, userService)
	st.NewHTTP(jwtBlocked, serverService)

	// Start the server
	e.Logger.Fatal(e.Start(":8090"))

	return nil
}
