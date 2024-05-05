package main

import (
	"encoding/json"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/util"
	"log"
)

func main() {
	es := util.NewElasticsearch()

	res, err := es.Client.Info()
	if err != nil {
		log.Fatalf("Error getting response from Elasticsearch: %s", err)
	}
	defer res.Body.Close()

	// Decode and print the cluster info
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	log.Printf("Elasticsearch cluster info: %v", r)
}
