package main

import (
	"log"

	"github.com/joho/godotenv"
	_ "github.com/mxngocqb/VCS-SERVER/back-end/docs"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
)

// @title Server Management System
// @version 1.0
// @host localhost:8090
// @BasePath /api
// @schemes http https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

func main() {
	_ = godotenv.Load()

	// load config
	cfgPath := "./conf.yaml"
	cfg, err := config.Load(cfgPath)

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// start API service
	if err := internal.Start(cfg); err != nil {
		log.Fatalf("Failed to start API: %v", err)
	}
}
