package health_check

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	service "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/server_status"
)

// StartPing function takes a serverService as input and starts the pingServer function
func pingServer(serverMap map[uint]service.Server, serverService *service.Service) {

	resultCh := make(chan []string)

	for _, server := range serverMap {
		go func(server service.Server) {
			ip := server.IP
			id := fmt.Sprintf("%d", server.ID)

			err := pingHost(ip, 50)
			status := "false"
			if err == nil {
				status = "true"
			}

			resultCh <- []string{id, status, ip}
		}(server)
	}

	for range serverMap {
		result := <-resultCh
		log.Printf("ID: %s, Server %s is up: %s\n", result[0], result[2], result[1])
		if result[1] == "true" {
			serverService.Update(result[0], true)
		} else {
			serverService.Update(result[0], false)
		}
	}
}

func getAllServer(servers *map[uint]service.Server, serverService *service.Service) {
	list, err := serverService.GetTotalServer()

	if err != nil {
		log.Fatalf("Failed to get total server: %v", err)
	}
	for _, server := range *list {
		if !server.DeletedAt.Valid {
			var serverForMap service.Server 
			serverForMap.ID = server.ID
			serverForMap.IP = server.IP
			serverForMap.Status = server.Status	
			(*servers)[server.ID] = serverForMap
		}
	}
}

// pingHost takes an IP address and a count of pings to send, and returns an error if any occurs.
func pingHost(host string, count int) error {
	cmd := exec.Command("ping", "-c", strconv.Itoa(count), "-W", "5", host)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// fetchServers function takes a URL and a JWT token as input and returns a ServersResponse and an error if any occurs
func fetchServers(url, token string) (*service.ServersResponse, error) {
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

	var serversResponse service.ServersResponse
	if err := json.NewDecoder(response.Body).Decode(&serversResponse); err != nil {
		return nil, err
	}

	return &serversResponse, nil
}

// Fetch, ping, and update server status
func fetchAndPingServers(url, token string, serverService *service.Service) {
	res, err := fetchServers(url, token)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	resultCh := make(chan []string)

	for _, server := range res.Data {
		go func(server service.Server) {
			ip := server.IP
			id := fmt.Sprintf("%d", server.ID)

			err := pingHost(ip, 50)
			status := "false"
			if err == nil {
				status = "true"
			}

			resultCh <- []string{id, status, ip}
		}(server)
	}

	for range res.Data {
		result := <-resultCh
		log.Printf("ID: %s, Server %s is up: %s\n", result[0], result[2], result[1])
		if result[1] == "true" {
			serverService.Update(result[0], true)
		} else {
			serverService.Update(result[0], false)
		}
	}
}
