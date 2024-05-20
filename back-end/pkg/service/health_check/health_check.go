package health_check

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	service "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/server_status"
)

func Start(cfg *config.Config) error{
	serverService, consumerService, err := Config(cfg)
	
	if err != nil {
		log.Fatalf("Failed to start server service: %v", err)
		return err
	}
	
	// Handle OS signals for graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	// Start the consumer
	serverMap := make(map[uint]service.Server)
	// Get all server from database
	getAllServer(&serverMap, serverService)
	log.Println("All server loaded")
	// Start the consumer to listen to messages from Kafka
	go consumerService.ConsumerStart(&serverMap, sigchan)
	log.Println("Consumer started, waiting for messages...")
	// Start the cron job
	StartPing(serverMap, serverService)
	// Keep the main program running until a signal is received
	<-sigchan

	return nil
}