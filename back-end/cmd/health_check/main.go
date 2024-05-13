package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-ping/ping"
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
	// Gửi yêu cầu GET đến API
	url := "http://localhost:8090/api/servers?limit=1000&offset=0"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJhZG1pbiIsImV4cCI6MTcxNTg3ODYzNX0.yLjMN8W6t6CP5Ghd3HHyebsNuhM4JR_OfgzH9iqUz6g" // Example JWT token
	
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer " + token)

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


func pingHost(ipAddr string) error {
	// Create a new Pinger
	pinger, err := ping.NewPinger(ipAddr)
	if err != nil {
		return err
	}

	// Set options for the pinger
	pinger.Count = 5 // Ping 5 times
	pinger.Timeout = time.Second * 5
	pinger.Interval = time.Second
	pinger.SetPrivileged(true) // This may be required in some operating systems to send ICMP packets

	// Variables to track successful pings
	successfulPings := 0

	// Start the ping loop
	pinger.OnRecv = func(pkt *ping.Packet) {
		successfulPings++
	}

	pinger.Run() // Blocks until finished

	// Check if all pings failed
	if successfulPings == 0 {
		return fmt.Errorf("no successful pings to %s", ipAddr)
	}

	return nil
}
