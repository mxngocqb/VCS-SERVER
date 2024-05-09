package main

import (
	"github.com/joho/godotenv"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	_ "github.com/mxngocqb/VCS-SERVER/back-end/docs"
	"log"
)

// @title VCS SMS API
// @description This is the server for the VCS SMS management system.
// @version 1.0

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
func main() {
	_ = godotenv.Load()

	cfgPath := "./conf.local.yaml"
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	if err := internal.Start(cfg); err != nil {
		log.Fatalf("failed to start API: %v", err)
	}
}
