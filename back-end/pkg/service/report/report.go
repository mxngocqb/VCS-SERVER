package report

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/server_status"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
)

func ScheduleDailyReport() {
	c := cron.New()
	// Send daily report at 8:00 AM
	_, err := c.AddFunc("8 9 * * *", func() {
	// _, err := c.AddFunc("@every 10m", func() {
		now := time.Now()
		loc, _ := time.LoadLocation("Asia/Bangkok") // Ensure timezone consistency with server logs
		start := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, loc)
		end := start.AddDate(0, 0, 1)
		err1 := SendReport([]string{"mxn111333@gmail.com"}, start, end)
		if err1 != nil {
			log.Printf("Error sending daily report: %v", err1)
		} else {
			log.Printf("Daily report sent successfully")
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling daily report: %v", err)
	}
	c.Start()
}

func FetchServersInfo(start, end time.Time) (float64, int, int, int, error) {
	es := service.NewElasticsearch() // assuming util contains the Elasticsearch service initialization

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
					Key int `json:"key"`
					// Key string `json:"key"`
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
        }`, /*bucket.Key*/ strconv.Itoa(bucket.Key))

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
		lastStatusMap[ /*bucket.Key*/ strconv.Itoa(bucket.Key)] = lastStatus
		if lastStatus {
			onlineServers++
		}

		// uptime, err := es.CalculateServerUptime( /*bucket.Key*/ strconv.Itoa(bucket.Key), now)
		// if err != nil {
		// 	log.Printf("Error calculating uptime for server %s: %v" /*strconv.Itoa(bucket.Key)*/, bucket.Key, err)
		// 	continue
		// }

		uptime, err := es.CalculateServerUptimeFromStartToEnd( /*bucket.Key*/ strconv.Itoa(bucket.Key), start, end)
		if err != nil {
			log.Printf("Error calculating uptime for server %s: %v" /*strconv.Itoa(bucket.Key)*/, bucket.Key, err)
			continue
		}

		totalUptime += uptime
	}

	if totalServers == 0 {
		return 0, 0, 0, 0, fmt.Errorf("no server data found for today")
	}

	if onlineServers == 0 {
		return 0, 0, totalServers, totalServers, nil
	}

	avgUptime := totalUptime.Hours() / float64(onlineServers)
	offlineServers := totalServers - onlineServers

	return avgUptime, onlineServers, offlineServers, totalServers, nil
}

// SendReport sends a daily server report to the specified email addresses.
// It fetches server information for the given time range and constructs an email
// containing the average uptime, number of online and offline servers, and the total number of servers.
// The email is sent using the SMTP protocol with the provided Gmail account credentials.
// If any error occurs during the process, it is returned. Otherwise, nil is returned.
func SendReport(email []string, start, end time.Time) error {
	avgUptime, online, offline, totalServers, err := FetchServersInfo(start, end)
	if err != nil {
		log.Printf("Error fetching server info: %v", err)
		return err
	}

	fmt.Println("Email:", email)

	m := gomail.NewMessage()
	m.SetHeader("From", "mxngocqb@gmail.com")
	m.SetHeader("To", email...)
	m.SetHeader("Subject", "Daily Server Report")
	m.SetBody("text/html", fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body {
				font-family: Arial, sans-serif;
				line-height: 1.6;
				color: #333;
			}
			.container {
				margin: 20px;
				padding: 20px;
				border: 1px solid #ddd;
				border-radius: 5px;
				background-color: #f9f9f9;
			}
			h2 {
				color: #4CAF50;
			}
			.data {
				margin-top: 20px;
			}
			.data strong {
				color: #555;
			}
			.footer {
				margin-top: 30px;
				font-size: 0.9em;
				color: #777;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>Daily Server Report</h2>
			<p>Hello,</p>
			<p>Here is the daily server report for your servers:</p>
			<p><strong>From:</strong> %s<br>
			<strong>To:</strong> %s</p>
			<div class="data">
				<p><strong>Average Uptime:</strong> %.2f hours</p>
				<p><strong>Online Servers:</strong> %d</p>
				<p><strong>Offline Servers:</strong> %d</p>
				<p><strong>Total Servers:</strong> %d</p>
			</div>
			<div class="footer">
				<p>Best regards,<br>
				Server Management Systems - VCS Team</p>
			</div>
		</div>
	</body>
	</html>`,
		start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"), avgUptime, online, offline, totalServers))

	d := gomail.NewDialer("smtp.gmail.com", 587, "mxngocqb@gmail.com", "xftw lchz hruo ojkq")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	} else {
		log.Println("Email sent successfully")
		return nil
	}
}
