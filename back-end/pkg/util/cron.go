package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
	"log"
	"strings"
	"time"
)

func ScheduleDailyReport() {
	c := cron.New()
	_, err := c.AddFunc("0 0 0 * * *", func() {
		now := time.Now()
		loc, _ := time.LoadLocation("Asia/Bangkok") // Ensure timezone consistency with server logs
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		end := start.AddDate(0, 0, 1)
		SendReport([]string{""}, start, end)
	})
	if err != nil {
		log.Fatalf("Error scheduling daily report: %v", err)
	}
	c.Start()
}

func FetchServersInfo(start, end time.Time) (float64, int, int, int, error) {
	es := NewElasticsearch() // assuming util contains the Elasticsearch service initialization

	// Today's date
	now := time.Now()

	fmt.Println("Start:", start)
	fmt.Println("End:", end)

	// Query to find all unique server IDs with logs today
	uniqueServersQuery := fmt.Sprintf(`
    {
        "size": 0,
        "aggs": {
            "unique_servers": {
                "terms": {
                    "field": "server_id",
                    "size": 10000 // Adjust size based on expected number of servers
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
					Key string `json:"key"`
				} `json:"buckets"`
			} `json:"unique_servers"`
		} `json:"aggregations"`
	}

	if err := json.NewDecoder(res.Body).Decode(&uniqueServersResp); err != nil {
		return 0, 0, 0, 0, err
	}

	totalServers := len(uniqueServersResp.Aggregations.UniqueServers.Buckets)
	totalUptime := time.Duration(0)
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
        }`, bucket.Key)

		lastStatusReq := esapi.SearchRequest{
			Index: []string{"server_status_logs"},
			Body:  strings.NewReader(lastStatusQuery),
		}

		lastStatusRes, err := lastStatusReq.Do(context.Background(), es.Client)
		if err != nil {
			log.Printf("Error fetching last status for server %s: %v", bucket.Key, err)
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
			log.Printf("Error decoding last status for server %s: %v", bucket.Key, err)
			continue
		}

		lastStatus := lastStatusResp.Hits.Hits[0].Source.Status
		lastStatusMap[bucket.Key] = lastStatus
		if lastStatus {
			onlineServers++
		}

		uptime, err := es.CalculateServerUptime(bucket.Key, now)
		if err != nil {
			log.Printf("Error calculating uptime for server %s: %v", bucket.Key, err)
			continue
		}
		totalUptime += uptime
	}

	if totalServers == 0 {
		return 0, 0, 0, 0, fmt.Errorf("no server data found for today")
	}

	avgUptime := totalUptime.Hours() / float64(onlineServers)
	offlineServers := totalServers - onlineServers

	return avgUptime, onlineServers, offlineServers, totalServers, nil
}

func SendReport(email []string, start, end time.Time) error {
	avgUptime, online, offline, totalServers, err := FetchServersInfo(start, end)
	if err != nil {
		return err
	}

	fmt.Println("Email:", email)

	m := gomail.NewMessage()
	m.SetHeader("From", "tuanthanh.vcs@gmail.com")
	m.SetHeader("To", email...)
	m.SetHeader("Subject", "Daily Server Report")
	m.SetBody("text/html", fmt.Sprintf(
		"<strong>Average Uptime:</strong> %.2f hours<br>"+
			"<strong>Online Servers:</strong> %d<br>"+
			"<strong>Offline Servers:</strong> %d<br>"+
			"<strong>Total Servers:</strong> %d",
		avgUptime, online, offline, totalServers))

	d := gomail.NewDialer("smtp.gmail.com", 587, "tuanthanh.vcs@gmail.com", "nwjs xikc oxya ygwc")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	} else {
		log.Println("Email sent successfully")
		return nil
	}
}
