package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	service "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/server_status"
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

	serverService, consumerService, err := Config(cfg)
	
	if err != nil {
		log.Fatalf("Failed to start server service: %v", err)
		return
	} 
	// Handle OS signals for graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	// Start the consumer
	serverMap := make(map[int]service.Server)
	go consumerService.ConsumerStart(&serverMap, sigchan)
	log.Println("Consumer started, waiting for messages...")
	// Start the cron job
	StartPing(serverMap, serverService)
    // Keep the main program running until a signal is received
    <-sigchan

}