package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"gorm.io/gorm"
)

// setupElasticsearchClient initializes and returns an Elasticsearch client configured for your environment.
func setupElasticsearchClient(t *testing.T) *ElasticService {
	elasticService := NewElasticsearch()
	if err := elasticService.CreateStatusLogIndex(); err != nil {
		t.Fatalf("Error creating status log index: %s", err)
	}
	return elasticService
}

func TestElasticService_IndexServer(t *testing.T) {
	es := setupElasticsearchClient(t)

	server := model.Server{
		Model: gorm.Model{
			ID:        1000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:   "localhost server",
		Status: true,
		IP:     "192.168.88.1",
	}

	err := es.IndexServer(server)
	if err != nil {
		t.Errorf("Error indexing server: %s", err)
	}

	t.Logf("Server indexed successfully")
	defer func(es *ElasticService, id string) {
		err := es.DeleteServerFromIndex(id)
		if err != nil {
			t.Errorf("Error deleting server: %s", err)
		}
	}(es, "1000")
}

func TestElasticService_DeleteServerFromIndex(t *testing.T) {
	es := setupElasticsearchClient(t)

	// First, index a server to ensure there is a document to delete
	server := model.Server{
		Model: gorm.Model{
			ID:        1001,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:   "localhost server",
		Status: true,
		IP:     "192.168.88.1",
	}

	// Index the server
	err := es.IndexServer(server)
	if err != nil {
		t.Fatalf("Error indexing server for deletion: %s", err)
	}

	// Wait a bit to ensure data consistency in Elasticsearch
	time.Sleep(1 * time.Second)

	// Now attempt to delete the server
	err = es.DeleteServerFromIndex("1001")
	if err != nil {
		t.Errorf("Error deleting server from index: %s", err)
	} else {
		t.Logf("Server deleted successfully")
	}
}

func TestElasticService_LogStatusChange(t *testing.T) {
	es := setupElasticsearchClient(t)

	// Define a server model
	server := model.Server{
		Model: gorm.Model{
			ID:        1002,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:   "localhost server",
		Status: true,
		IP:     "192.168.88.1",
	}

	// Log status change
	err := es.LogStatusChange(server, false)
	if err != nil {
		t.Errorf("Error logging status change: %s", err)
	} else {
		t.Logf("Status change logged successfully")
	}

	// Cleanup: Delete the log entry to maintain a clean state.
	// This might involve a direct delete by query or similar depending on setup.
	defer func(es *ElasticService, serverID uint) {
		err := es.DeleteServerLogs(strconv.Itoa(int(serverID)))
		if err != nil {
			t.Errorf("Error cleaning up server logs: %s", err)
		}
	}(es, server.ID)
}

func TestElasticService_CreateStatusLogIndex(t *testing.T) {
	es := setupElasticsearchClient(t)

	// Attempt to create the index
	err := es.CreateStatusLogIndex()
	if err != nil {
		t.Errorf("Error creating status log index: %s", err)
	} else {
		t.Logf("Status log index created successfully")
	}

	// Optionally, verify index settings and mappings if necessary
	// This would typically involve a separate call to check index configuration
	// This step is more relevant in integration testing

	// Cleanup: Consider deleting the index after testing to clean up the environment.
	// Note: Be very cautious with delete operations in a shared environment.
	defer func(es *ElasticService) {
		deleteIndexResponse, err := es.Client.Indices.Delete([]string{"server_status_logs"})
		if err != nil || deleteIndexResponse.IsError() {
			t.Logf("Failed to delete 'server_status_logs' index during cleanup: %s", err)
		}
	}(es)
}

func TestElasticService_DeleteServerLogs(t *testing.T) {
	es := setupElasticsearchClient(t)

	// Define a server model
	server := model.Server{
		Model: gorm.Model{
			ID:        1003, // Make sure to use a unique ID to avoid test data collisions
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:   "test server",
		Status: true,
		IP:     "123.45.67.89",
	}

	// Log multiple status changes for the server
	// For example, log status changes every hour for the last 6 hours
	for i := 0; i < 6; i++ {
		status := true // Toggle status every hour
		err := es.LogStatusChange(server, status)
		if err != nil {
			t.Fatalf("Error logging status change: %s", err)
		}
		time.Sleep(time.Second) // Sleep for a second to ensure unique timestamps
	}

	// Ensure logs are indexed
	time.Sleep(2 * time.Second) // Wait for logs to be indexed

	// Delete the log entries
	err := es.DeleteServerLogs(strconv.Itoa(int(server.ID)))
	if err != nil {
		t.Errorf("Error deleting server logs: %s", err)
	} else {
		t.Logf("Server logs deleted successfully")
	}

	// Verify that the logs are deleted
	query := fmt.Sprintf(`
    {
        "query": {
            "term": {
                "server_id": "%d"
            }
        }
    }`, server.ID)

	req := esapi.SearchRequest{
		Index: []string{"server_status_logs"},
		Body:  strings.NewReader(query),
	}
	// Ensure logs are indexed
	time.Sleep(2 * time.Second) // Wait for logs to be indexed
	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		t.Fatalf("Error searching for logs after deletion: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		t.Fatalf("Error response from Elasticsearch: %s", res.String())
	}

	// Parse the response to check the log count
	var r struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		t.Fatalf("Error parsing search response: %s", err)
	}

	if r.Hits.Total.Value != 0 {
		t.Errorf("Expected 0 logs, but found %d logs", r.Hits.Total.Value)
	} else {
		t.Logf("Verified that logs are deleted successfully")
	}
}
func TestElasticService_CalculateServerUptime(t *testing.T) {
	es := setupElasticsearchClient(t)

	serverID := "1004"
	date := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)

	uptime, err := es.CalculateServerUptime(serverID, date)
	if err != nil {
		t.Errorf("Error calculating server uptime: %s", err)
	}

	expectedUptime := time.Hour * 0
	if uptime != expectedUptime {
		t.Errorf("Expected uptime to be %s, but got %s", expectedUptime, uptime)
	} else {
		t.Logf("Server uptime calculated successfully")
	}
}

func TestElasticService_CalculateServerUptimeFromStartToEnd(t *testing.T) {
	es := setupElasticsearchClient(t)

	serverID := "1004"
	startDate := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.Add(time.Hour * 24)
	uptime, err := es.CalculateServerUptimeFromStartToEnd(serverID, startDate, endDate)
	if err != nil {
		t.Errorf("Error calculating server uptime: %s", err)
	}

	expectedUptime := time.Hour * 0
	if uptime != expectedUptime {
		t.Errorf("Expected uptime to be %s, but got %s", expectedUptime, uptime)
	} else {
		t.Logf("Server uptime calculated successfully")
	}
}