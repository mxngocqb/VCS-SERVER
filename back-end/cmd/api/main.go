package main

import (
	"github.com/joho/godotenv"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	_ "github.com/mxngocqb/VCS-SERVER/back-end/docs"
	"log"
)

// @title Your API Title
// @version 1.0
// @description This is a sample API
// @host localhost:8090
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

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
