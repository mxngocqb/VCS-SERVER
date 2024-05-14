package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-ping/ping"
	"github.com/joho/godotenv"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	service "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/server_status"
)

type Server struct {
	ID        int    `json:"ID"`
	CreatedAt string `json:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt"`
	DeletedAt string `json:"DeletedAt"`
	Name      string `json:"name"`
	Status    bool   `json:"status"`
	IP        string `json:"ip"`
}

type ServersResponse struct {
	Total int      `json:"total"`
	Data  []Server `json:"data"`
}

func main() {
	url := "http://localhost:8090/api/servers?limit=1000&offset=0"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJhZG1pbiIsImV4cCI6MTcxNTg3ODYzNX0.yLjMN8W6t6CP5Ghd3HHyebsNuhM4JR_OfgzH9iqUz6g" // Example JWT token

	// Load environment variables
	_ = godotenv.Load()

	// load config
	cfgPath := "./conf.yaml"
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := service.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repository := service.NewServerRepository(db.DB)
	elasticService := service.NewElasticsearch()
	serverService := service.NewServerService(repository, elasticService)

	// fetch data from the API
	res, err := fetchServers(url, token)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	resultCh := make(chan []string)

	// Loop through each server in the response data
	for _, server := range res.Data {
		go func(server Server) {
			ip := server.IP
			id := fmt.Sprintf("%d", server.ID)

			// Ping the IP address
			err := pingHost(ip, 10)
			status := "false"
			if err == nil {
				status = "true"
			}

			resultCh <- []string{id, status, ip}
		}(server)
	}

	// Collect the results from the channel
	for range res.Data {
		result := <-resultCh
		fmt.Printf("ID: %s, Status: %s, IP: %s\n", result[0], result[1], result[2]) 
		if result[1] == "true" {
			serverService.Update(result[0], true)
		} else {
			serverService.Update(result[0], false)
		}
	}
}

// ping function takes an IP address as input and returns an error if any occurs
func pingHost(host string, count int) error {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return err
	}

	pinger.Count = count
	pinger.Timeout = time.Second * time.Duration(count)
	pinger.SetPrivileged(true)

	pinger.OnRecv = func(pkt *ping.Packet) {
	}
	pinger.OnFinish = func(stats *ping.Statistics) {

	}

	err = pinger.Run()
	if err != nil {
		return err
	}

	stats := pinger.Statistics()
	success := stats.PacketsRecv > 0
	if !success {
		return fmt.Errorf("no packets received")
	}

	return nil
}

// fetchServers function takes a URL and a JWT token as input and returns a ServersResponse and an error if any occurs
func fetchServers(url, token string) (*ServersResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var serversResponse ServersResponse
	if err := json.NewDecoder(response.Body).Decode(&serversResponse); err != nil {
		return nil, err
	}

	return &serversResponse, nil
}
