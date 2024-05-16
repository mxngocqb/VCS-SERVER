package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"log"
	"strconv"
	"strings"
	"time"

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
	} else {
		log.Printf("Connected to Elasticsearch")
	}

	// Check cluster health to ensure proper connection
	res, err := es.Info()
	if err != nil || res.IsError() {
		log.Fatalf("Error connecting to Elasticsearch at startup: %s", err)
	}

	return &ElasticService{es}
}

// IndexServer indexes or updates a server document in Elasticsearch.
func (es *ElasticService) IndexServer(server model.Server) error {
	data, err := json.Marshal(server)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "servers",
		DocumentID: strconv.Itoa(int(server.ID)),
		Body:       strings.NewReader(string(data)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document ID %d: %s", server.ID, res.String())
	}

	return nil
}

// DeleteServerFromIndex removes a server document from Elasticsearch.
func (es *ElasticService) DeleteServerFromIndex(id string) error {
	req := esapi.DeleteRequest{
		Index:      "servers",
		DocumentID: id,
	}

	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting document ID %s: %s", id, res.String())
	}

	return nil
}

// LogStatusChange logs a server's status change to Elasticsearch.
func (es *ElasticService) LogStatusChange(server model.Server, status bool) error {
	logEntry := map[string]interface{}{
		"server_id": server.ID,
		"status":    status,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	data, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:   "server_status_logs",
		Body:    strings.NewReader(string(data)),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error logging status change for server ID %d: %s", server.ID, res.String())
	}

	return nil
}

// CalculateServerUptime calculates the total uptime for a server on a given day, accounting for timezone.
func (es *ElasticService) CalculateServerUptime(serverID string, date time.Time) (time.Duration, error) {
	var totalUptime time.Duration

	// Assume the server timestamp is in GMT+7
	loc, _ := time.LoadLocation("Asia/Bangkok") // Load the GMT+7 location

	// Define the start and end of the day in the server's timezone
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Nanosecond)
	now := time.Now().In(loc)

	// Elasticsearch query
	query := fmt.Sprintf(`
    {
        "query": {
            "bool": {
                "must": [
                    {"term": {"server_id": %s}},
                    {"range": {"timestamp": {"gte": "%s", "lte": "%s"}}}
                ]
            }
        },
        "sort": [{"timestamp": {"order": "asc"}}]
    }`, serverID, startOfDay.Format(time.RFC3339), endOfDay.Format(time.RFC3339))

	req := esapi.SearchRequest{
		Index: []string{"server_status_logs"},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("error searching logs for server ID %s: %s", serverID, res.String())
	}

	// Parse the response to calculate total uptime
	var r struct {
		Hits struct {
			Hits []struct {
				Source struct {
					Timestamp time.Time `json:"timestamp"`
					Status    bool      `json:"status"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return 0, fmt.Errorf("error parsing response: %s", err)
	}

	// Track the last "on" time; if the server never turns "off", it's on till the end of the day
	var lastOnTime *time.Time

	for _, hit := range r.Hits.Hits {
		if hit.Source.Status {
			// Server turned "on"
			if lastOnTime != nil {
				totalUptime += hit.Source.Timestamp.Sub(*lastOnTime)
			}
			lastOnTime = &hit.Source.Timestamp
		} else if lastOnTime != nil {
			// Server turned "off"
			totalUptime += hit.Source.Timestamp.Sub(*lastOnTime)
			lastOnTime = nil // Reset lastOnTime after calculating uptime
		}
	}

	// If the last status was "on" and there was no "off" event, count uptime till end of the day
	if lastOnTime != nil && endOfDay.Before(now) {
		totalUptime += endOfDay.Sub(*lastOnTime)
	} else if lastOnTime != nil && endOfDay.After(now) {
		totalUptime += now.Sub(*lastOnTime)
	}

	return totalUptime, nil
}

// CreateStatusLogIndex creates an index for server documents in Elasticsearch.
func (es *ElasticService) CreateStatusLogIndex() error {
	// Check if the index already exists
	res, err := es.Client.Indices.Exists([]string{"server_status_logs"})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// If the index doesn't exist, create it
	if res.StatusCode == 404 {
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

// DeleteServerLogs removes all log entries for a specific server from Elasticsearch.
func (es *ElasticService) DeleteServerLogs(serverID string) error {
	// Elasticsearch query to delete documents
	query := fmt.Sprintf(`
    {
        "query": {
            "term": {
                "server_id": %s
            }
        }
    }`, serverID)

	req := esapi.DeleteByQueryRequest{
		Index: []string{"server_status_logs"},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return fmt.Errorf("error sending delete request for server ID %s: %s", serverID, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting logs for server ID %s: %s", serverID, res.String())
	}

	return nil
}
