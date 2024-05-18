package main

import (
	"log"
	"github.com/joho/godotenv"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/health_check"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// load config
	cfgPath := "./conf.yaml"
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	err = health_check.Start(cfg)
	if err != nil {
		log.Fatalf("Error starting server service: %v", err)
	}
	

}
