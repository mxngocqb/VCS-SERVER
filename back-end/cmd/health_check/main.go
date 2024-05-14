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
	service.NewServerService(repository, elasticService)

	// Gửi yêu cầu GET đến API
	url := "http://localhost:8090/api/servers?limit=1000&offset=0"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJhZG1pbiIsImV4cCI6MTcxNTkxMTY1OH0.s3l626nOilzilRDfcsOS5tFdqc5oQqmqTUEgjrTNv9k" // Example JWT token

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	response, err := client.Do(req)

	if err != nil {
		log.Fatal("Lỗi khi gửi yêu cầu:", err)
	}
	defer response.Body.Close()

	// Đảm bảo rằng mã trạng thái là 200 OK
	if response.StatusCode != http.StatusOK {
		log.Fatalf("Lỗi: mã trạng thái không phải là 200 OK. Mã: %d", response.StatusCode)
	}

	// Phân tích dữ liệu JSON từ phản hồi
	var serversResponse ServersResponse
	if err := json.NewDecoder(response.Body).Decode(&serversResponse); err != nil {
		log.Fatal("Lỗi khi phân tích dữ liệu JSON:", err)
	}

	// Trích xuất mảng IP từ dữ liệu
	var ips []string
	for _, server := range serversResponse.Data {
		ips = append(ips, server.IP)
	}

	resultCh := make(chan []string)

	// Loop through each record in the input CSV file
	for _, record := range ips {
		go func(record string) {
			ip := record

			// Ping the ip address
			err := pingHost(ip)
			status := "false"
			if err == nil {
				status = "true"
			}
			resultCh <- []string{ip, status}
		}(record)
	}

	// Collect the results from the channel
	for range ips {
		result := <-resultCh
		fmt.Printf("IP: %s, Status: %s\n", result[0], result[1])
	}
}

// ping function takes an IP address as input and returns an error if any occurs
func pingHost(ip string) error {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return fmt.Errorf("error creating pinger: %v", err)
	}

	// Set options for the pinger
	pinger.Count = 3                 // Number of packets to send
	pinger.Interval = time.Second    // Time between each ping
	pinger.Timeout = time.Second * 5 // Timeout for the entire ping operation

	// Define a callback to handle each received packet (optional)
	pinger.OnRecv = func(pkt *ping.Packet) {}

	// Run the ping
	err = pinger.Run()
	if err != nil {
		return fmt.Errorf("error running pinger: %v", err)
	}

	return nil
}
