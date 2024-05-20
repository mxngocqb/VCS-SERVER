package elastic

import (
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
)

func ConnectElasticSearch(config *config.Config)( *elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: config.ELASTIC.Hosts,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
		return nil, err
	} else {
		log.Printf("Connected to Elasticsearch")
	}

	// Check cluster health to ensure proper connection
	// res, err := es.Info()
	// if err != nil || res.IsError() {
	// 	log.Fatalf("Error connecting to Elasticsearch at startup: %s", err)
	// 	return nil, err
	// }

	return es, err
}