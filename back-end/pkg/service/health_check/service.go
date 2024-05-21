package health_check

import (
	"log"

	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/elastic"
	kafka "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/kafka"
	service "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/server_status"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/util"
	"github.com/robfig/cron/v3"
)

// Config sets up the server service.
func Config(cfg *config.Config) (*service.Service, *kafka.ConsumerService, error) {
	logger, err := util.NewPostgresLogger()
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}
	db, err := service.New(cfg, logger)
	if err != nil {
		return nil, nil, err
	}

	repository := service.NewServerRepository(db.DB)
	elasticClient, err := elastic.ConnectElasticSearch(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to ElasticSearch: %v", err)
	}
	elasticService := elastic.NewElasticsearch(elasticClient)
	serverService := service.NewServerService(repository, elasticService)
	consumerService := kafka.NewConsumerService(cfg)
	

	return serverService, consumerService, nil
}

// Start starts the cron job.
func StartPing(serverMap map[uint]service.Server, serverService *service.Service){
	c := cron.New()

	_, err := c.AddFunc("@every 5m", func() {
		pingServer(serverMap, serverService)
	})

	if err != nil {
		log.Fatalf("Error scheduling daily report: %v", err)
	}

	c.Start()
}