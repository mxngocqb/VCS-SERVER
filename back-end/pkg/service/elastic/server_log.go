package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
)

type ElasticService interface {
	IndexServer(server model.Server) error
	DeleteServerFromIndex(id string) error
	LogStatusChange(server model.Server, status bool) error
	CalculateServerUptime(serverID string, date time.Time) (time.Duration, error)
	CalculateServerUptimeFromStartToEnd(serverID string, startDate, endDate time.Time) (time.Duration, error) 
	CreateStatusLogIndex() error
	DeleteServerLogs(serverID string) error
	FetchServersInfo(start, end time.Time) (float64, int, int, int, error)
}

// ElasticServiceImpl provides methods to interact with Elasticsearch.
type ElasticServiceImpl struct {
	Client *elasticsearch.Client
}

// NewElasticsearch initializes and returns an Elasticsearch client configured for your environment.
func NewElasticsearch(elasticClient *elasticsearch.Client) ElasticService {
		return &ElasticServiceImpl{elasticClient}
}

// IndexServer indexes or updates a server document in Elasticsearch.
func (es *ElasticServiceImpl) IndexServer(server model.Server) error {
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
func (es *ElasticServiceImpl) DeleteServerFromIndex(id string) error {
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
func (es *ElasticServiceImpl) LogStatusChange(server model.Server, status bool) error {
	var logEntry map[string]interface{}
	if status == true{
		logEntry = map[string]interface{}{
			"server_id": server.ID,
			"status":    status,
			"timestamp": time.Now().Format(time.RFC3339),
			"timeduration": 5,
		}
	} else{
		logEntry = map[string]interface{}{
			"server_id": server.ID,
			"status":    status,
			"timestamp": time.Now().Format(time.RFC3339),
			"timeduration": 0,
		}
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
func (es *ElasticServiceImpl) CalculateServerUptime(serverID string, date time.Time) (time.Duration, error) {
	var totalUptime time.Duration

	// Assume the server timestamp is in GMT+7
	loc, _ := time.LoadLocation("Asia/Bangkok") // Load the GMT+7 location

	// Define the start and end of the day in the server's timezone
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Nanosecond)
	// Elasticsearch query
	query := fmt.Sprintf(`
    {
		"size": 0,
        "query": {
            "bool": {
                "must": [
                    {"term": {"server_id": %s}},
                    {"range": {"timestamp": {"gte": "%s", "lte": "%s"}}}
                ]
            }
        },
		"aggs":{
			"uptime": {
				"sum": {
					"field": "timeduration"
				}
			}
		  }
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
		Aggregations struct {
			Uptime struct {
				Value float32 `json:"value"`
			} `json:"uptime"`
		} `json:"aggregations"`
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return 0, fmt.Errorf("error parsing response: %s", err)
	}

	totalUptime = time.Duration(r.Aggregations.Uptime.Value) * time.Minute
	
	return totalUptime , nil
}

// CalculateServerUptime calculates the total uptime for a server on a given day, accounting for timezone.
func (es *ElasticServiceImpl) CalculateServerUptimeFromStartToEnd(serverID string, startDate, endDate time.Time) (time.Duration, error) {
	var totalUptime time.Duration

	// Assume the server timestamp is in GMT+7
	loc, _ := time.LoadLocation("Asia/Bangkok") // Load the GMT+7 location

	// Define the start and end of the day in the server's timezone
	startOfDay := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
	endOfDay :=  time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, loc)
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
func (es *ElasticServiceImpl) CreateStatusLogIndex() error {
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
func (es *ElasticServiceImpl) DeleteServerLogs(serverID string) error {
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

// FetchServersInfo retrieves the average uptime, number of online servers, and total servers for a given time range.
func (es *ElasticServiceImpl) FetchServersInfo(start, end time.Time) (float64, int, int, int, error) {

	fmt.Println("Start:", start)
	fmt.Println("End:", end)

	// Elasticsearch query to get unique servers
	uniqueServersQuery := fmt.Sprintf(`
	{
        "size": 0,
        "aggs": {	
            "unique_servers": {
                "terms": {
                    "field": "server_id",
                    "size": 10000
                },
				"aggs": {
					"total_uptime": {
						"sum": {
							"field": "timeduration"
						}
					}
				}
            },
			"average_uptime_per_server": {
				"avg_bucket": {
					"buckets_path": "unique_servers>total_uptime"
				}
			}
        },
        "query": {
            "range": {
                "timestamp": {
                    "gte": "%s",
                    "lte": "%s"
                }
            }
        }
    }`, start.Format(time.RFC3339), end.Format(time.RFC3339))

	req := esapi.SearchRequest{
		Index: []string{"server_status_logs"},
		Body:  strings.NewReader(uniqueServersQuery),
	}

	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	defer res.Body.Close()

	// Parse unique servers response
	var uniqueServersResp struct {
		Aggregations struct {
			UniqueServers struct {
				Buckets []struct {
					// Key int `json:"key"`
					Key string `json:"key"`
				} `json:"buckets"`
			} `json:"unique_servers"`
			AverageUptimePerServer struct {
				Value float64 `json:"value"`
			} `json:"average_uptime_per_server"`
		} `json:"aggregations"`
	}

	if err := json.NewDecoder(res.Body).Decode(&uniqueServersResp); err != nil {
		return 0, 0, 0, 0, err
	}

	totalServers := len(uniqueServersResp.Aggregations.UniqueServers.Buckets)
	avgUptime := uniqueServersResp.Aggregations.AverageUptimePerServer.Value/60
	onlineServers := 0

	// Store the last status of each server
	lastStatusMap := make(map[string]bool)

	// Additional Query to get the last status for each server
	for _, bucket := range uniqueServersResp.Aggregations.UniqueServers.Buckets {
		lastStatusQuery := fmt.Sprintf(`
        {
            "query": {
                "term": {
                    "server_id": "%s"
                }
            },
            "size": 1,
            "sort": [
                {
                    "timestamp": {
                        "order": "desc"
                    }
                }
            ]
        }`, /*strconv.Itoa(bucket.Key)*/ bucket.Key)

		lastStatusReq := esapi.SearchRequest{
			Index: []string{"server_status_logs"},
			Body:  strings.NewReader(lastStatusQuery),
		}

		lastStatusRes, err := lastStatusReq.Do(context.Background(), es.Client)
		if err != nil {
			log.Printf("Error fetching last status for server %s: %v" /*strconv.Itoa(bucket.Key)*/, bucket.Key, err)
			continue
		}
		defer lastStatusRes.Body.Close()

		var lastStatusResp struct {
			Hits struct {
				Hits []struct {
					Source struct {
						Status bool `json:"status"`
					} `json:"_source"`
				} `json:"hits"`
			} `json:"hits"`
		}

		if err := json.NewDecoder(lastStatusRes.Body).Decode(&lastStatusResp); err != nil {
			log.Printf("Error decoding last status for server %s: %v" /*strconv.Itoa(bucket.Key)*/, bucket.Key, err)
			continue
		}

		lastStatus := lastStatusResp.Hits.Hits[0].Source.Status
		lastStatusMap[ /* strconv.Itoa(bucket.Key)*/ bucket.Key] = lastStatus
		if lastStatus {
			onlineServers++
		}
	}
	

	if totalServers == 0 {
		return 0, 0, 0, 0, fmt.Errorf("no server data found for today")
	}

	if onlineServers == 0 {
		return 0, 0, totalServers, totalServers, nil
	}

	// avgUptime := totalUptime.Hours() / float64(onlineServers)
	offlineServers := totalServers - onlineServers

	return avgUptime, onlineServers, offlineServers, totalServers, nil
}