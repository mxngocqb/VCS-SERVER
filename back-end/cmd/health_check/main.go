package main

import (
	"log"
	"github.com/joho/godotenv"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
)

func main() {
	url := "http://localhost:8090/api/servers?limit=1000&offset=0"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJhZG1pbiIsImV4cCI6MTcxNTg3ODYzNX0.yLjMN8W6t6CP5Ghd3HHyebsNuhM4JR_OfgzH9iqUz6g" // Example JWT token

	// Load environment variables
	_ = godotenv.Load()

	// load config
	cfgPath := "./conf.yaml"
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	serverService, err := Config(cfg)
	
	if err != nil {
		log.Fatalf("Failed to start server service: %v", err)
		return
	} 

	Start(url, token, serverService)

	select {}
}