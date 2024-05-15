package main

import (
	"log"

	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	service "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/server_status"
	"github.com/robfig/cron/v3"
)

// Config sets up the server service.
func Config(cfg *config.Config) (*service.Service, error) {

	db, err := service.New(cfg)
	if err != nil {
		return nil, err
	}

	repository := service.NewServerRepository(db.DB)
	elasticService := service.NewElasticsearch()
	serverService := service.NewServerService(repository, elasticService)

	return serverService, nil
}

// Start starts the cron job.
func Start(url, token string, serverService *service.Service){
	c := cron.New()

	_, err := c.AddFunc("@every 5m", func() {
		fetchAndPingServers(url, token, serverService)
	})

	if err != nil {
		log.Fatalf("Error scheduling daily report: %v", err)
	}

	c.Start()
}