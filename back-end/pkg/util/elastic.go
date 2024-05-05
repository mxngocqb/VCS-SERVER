package util

import (
	"log"
	"strings"
	"fmt"
	// "github.com/elastic/go-elasticsearch/v8/esapi"
	// "github.com/mxngocqb/VCS-SERVER/back-end/internal/model"

	"github.com/elastic/go-elasticsearch/v7"
)

// ElasticService provides methods to interact with Elasticsearch.
type ElasticService struct {
	Client *elasticsearch.Client
}

// NewElasticsearch initializes and returns an Elasticsearch client configured for your environment.
func NewElasticsearch() *ElasticService {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200", // Elasticsearch URL
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
	} else{
		log.Printf("Connected to Elasticsearch")
	}

	// Check cluster health to ensure proper connection
	res, err := es.Info()
	if err != nil || res.IsError() {
		log.Fatalf("Error connecting to Elasticsearch at startup: %s", err)
	}

	return &ElasticService{es}
}

// CreateStatusLogIndex creates an index for server documents in Elasticsearch.
func (es *ElasticService) CreateStatusLogIndex() error {
	// Check if the index already exists
	existsRes, err := es.Client.Indices.Exists([]string{"server_status_logs"})
	if err != nil {
		return err
	}
	defer existsRes.Body.Close()

	// If the index doesn't exist, create it
	if existsRes.StatusCode == 404 {
		createRes, err := es.Client.Indices.Create(
			"server_status_logs",
			es.Client.Indices.Create.WithBody(strings.NewReader(`{
                "mappings": {
                    "properties": {
                        "server_id": { "type": "keyword" },
                        "status": { "type": "boolean" },
                        "timestamp": { "type": "date" }
                    }
                }
            }`)),
		)
		if err != nil {
			return err
		}
		defer createRes.Body.Close()

		if createRes.IsError() {
			return fmt.Errorf("error creating index: %s", createRes.String())
		}
	}

	return nil
}