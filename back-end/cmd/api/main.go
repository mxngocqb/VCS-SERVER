package main

import (
	"github.com/joho/godotenv"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	_ "github.com/mxngocqb/VCS-SERVER/back-end/docs"
	"log"
)

// @title Viettel Cyber Security - Server Management System
// @version 1.0
// @description This is the API documentation for the Viettel Cyber Security - Server Management System.
// @host 192.168.88.130:8090
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
