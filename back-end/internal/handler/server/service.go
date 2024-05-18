package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	// labstak/echo is a web framework for Go
	"github.com/labstack/echo/v4"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/handler"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/repository"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/cache"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/kafka"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/elastic"
	pb "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/report/proto"
	util "github.com/mxngocqb/VCS-SERVER/back-end/pkg/util"
	"gorm.io/gorm"

	// gRPC framework for Go
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IServerService interface {
	View(c echo.Context, perPage int, offset int, status, field, order string) ([]model.Server, int64, error)
	Create(c echo.Context, server *model.Server) (*model.Server, error)
	CreateMany(c echo.Context, servers []model.Server) ([]model.Server, []int, []int, error)
	Update(c echo.Context, id string, server *model.Server) (*model.Server, error)
	Delete(c echo.Context, id string) error
	GetServersFiltered(c echo.Context, startCreated, endCreated, startUpdated, endUpdated, field, order string) error
	GetServerUptime(c echo.Context, serverID string, date string) (time.Duration, error)
	GetServerReport(c echo.Context, mail, start, end string) error
}

type Service struct {
	repository *repository.ServerRepository
	rbac       *handler.RbacService
	elastic    elastic.ElasticService
	cache      cache.ServerCache
	producer   *kafka.ProducerService
}

func NewServerService(repository *repository.ServerRepository, rbac *handler.RbacService, elastic elastic.ElasticService, sc cache.ServerCache,
	producer *kafka.ProducerService) *Service {
	return &Service{
		repository: repository,
		rbac:       rbac,
		elastic:    elastic,
		cache:      sc,
		producer:   producer,
	}
}

// View retrieves servers from the database with optional pagination and status filtering.
func (s *Service) View(c echo.Context, perPage int, offset int, status, field, order string) ([]model.Server, int64, error) {
	// ctx := c.Request().Context()
	key := s.cache.ConstructCacheKey(perPage, offset, status, field, order)
	key_total := key + "_total"
	// Try to get data from Redis first
	values := s.cache.GetMultiRequest(key)
	total := s.cache.GetTotalServer(key_total)

	if values != nil && total != -1 {
		// Data found in cache
		return values, total, nil
	}

	// Data not found in cache, fetch from database
	servers, numberOfServer, err := s.repository.GetServersFiltered(perPage, offset, status, field, order)
	if err != nil {
		return nil, -1, err
	}

	s.cache.SetMultiRequest(key, servers) // Adjust expiration as needed
	s.cache.SetTotalServer(key_total, numberOfServer)
	return servers, numberOfServer, nil
}

// Create creates a new server.
func (s *Service) Create(c echo.Context, server *model.Server) (*model.Server, error) {
	s.cache.InvalidateCache()
	fmt.Println("Create call")
	// Role ID for creating a new server is 1 (Admin)
	requiredRoleID := uint(1)

	// Enforce role check
	if err := s.rbac.EnforceRole(c, requiredRoleID); err != nil {
		return &model.Server{}, err // This will handle forbidden access
	}

	// Create new server in the database
	err := s.repository.Create(server)
	if err != nil {
		return &model.Server{}, err
	}

	fmt.Println("server", server)

	err = s.elastic.IndexServer(*server)
	if err != nil {
		return &model.Server{}, err
	}

	// After successfully creating the server, log the status change
	err = s.elastic.LogStatusChange(*server, server.Status)
	if err != nil {
		// Handle logging error, you may choose to return an error or just log it
		return &model.Server{}, err
	}

	s.producer.SendServer(server.ID, *server)

	// Cache the server
	s.cache.Set(strconv.Itoa(int(server.ID)), server)

	return server, nil
}

// CreateMany creates multiple servers and returns detailed results.
func (s *Service) CreateMany(c echo.Context, servers []model.Server) ([]model.Server, []int, []int, error) {
	s.cache.InvalidateCache()
	requiredRoleID := uint(1)
	if err := s.rbac.EnforceRole(c, requiredRoleID); err != nil {
		return nil, nil, nil, err
	}

	var createdServers []model.Server
	var successLines, failedLines []int

	for i, server := range servers {
		err := s.repository.Create(&server)
		if err != nil {
			failedLines = append(failedLines, i+2) // +2 to account for zero index and header row
			continue
		}
		createdServers = append(createdServers, server)
		successLines = append(successLines, i+2)

		err = s.elastic.IndexServer(server)
		if err != nil {
			log.Printf("Error indexing server %d: %v", server.ID, err)
		}
		// After successfully creating the server, log the status change
		err = s.elastic.LogStatusChange(server, server.Status)
		s.producer.SendServer(server.ID, server)
		if err != nil {
			// Handle logging error
			log.Printf("Error logging status change for server ID in Elasticsearch", server.ID)
		}
	}

	return createdServers, successLines, failedLines, nil
}

// Update updates a server.
func (s *Service) Update(c echo.Context, id string, server *model.Server) (*model.Server, error) {
	s.cache.InvalidateCache()
	// Role ID for updating a server is 1 (Admin)
	requiredRoleID := uint(1)

	// Enforce role check
	if err := s.rbac.EnforceRole(c, requiredRoleID); err != nil {
		return &model.Server{}, err // This will handle forbidden access
	}

	// Update server in Elasticsearch
	existingServer, err := s.repository.GetServerByID(id)
	if err != nil {
		return nil, err
	}

	existingServerStatus := existingServer.Status

	// Update server in the database
	err = s.repository.Update(id, server)
	if err != nil {
		return &model.Server{}, err
	}
	// Retrieve updated server
	updatedServer, err := s.repository.GetServerByID(id)
	if err != nil {
		return &model.Server{}, err
	}

	if existingServerStatus != updatedServer.Status {
		err = s.elastic.LogStatusChange(*updatedServer, server.Status)
		if err != nil {
			return nil, err
		}
	}

	// Cache the server
	s.producer.SendServer(server.ID, *server)
	s.cache.Set(strconv.Itoa(int(server.ID)), server)

	return updatedServer, nil
}

// Delete deletes a server.
func (s *Service) Delete(c echo.Context, id string) error {
	s.cache.InvalidateCache()
	// Role ID for deleting a server is 1 (Admin)
	requiredRoleID := uint(1)

	// Enforce role check
	if err := s.rbac.EnforceRole(c, requiredRoleID); err != nil {
		return err // This will handle forbidden access
	}

	server, err := s.repository.GetServerByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Server with ID %s not found", id))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve server: %v", err))
	}

	// Log status change before deleting the server
	err = s.elastic.LogStatusChange(*server, false)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error logging status change: %v", err))
	}

	// Delete server from the database
	err = s.repository.Delete(id)
	if err != nil {
		return err
	}

	// DELETE FROM ELASTICSEARCH MIGHT CAUSE ERROR IF THE SERVERS ARE NOT CREATED USING THE ENDPOINT (THEY ARE NOT CREATED IN ELASTICSEARCH IF USING SQL COMMAND ONLY)
	// Delete server from Elasticsearch
	err = s.elastic.DeleteServerFromIndex(id)
	if err != nil {
		return err
	}

	err = s.elastic.DeleteServerLogs(id)
	if err != nil {
		return err
	}
	// Cache the server
	s.cache.Delete(strconv.Itoa(int(server.ID)))
	s.producer.DropServer(server.ID)
	return nil
}

// GetServersFiltered retrieves servers with optional date range filtering.
func (s *Service) GetServersFiltered(c echo.Context, startCreated, endCreated, startUpdated, endUpdated, field, order string) error {
	layout := "2006-01-02"

	var startCreatedTime, endCreatedTime, startUpdatedTime, endUpdatedTime time.Time
	var err error

	if startCreated != "" && endCreated != "" {
		startCreatedTime, err = time.Parse(layout, startCreated)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid startCreated date format")
		}
		endCreatedTime, err = time.Parse(layout, endCreated)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid endCreated date format")
		}
	}

	if startUpdated != "" && endUpdated != "" {
		startUpdatedTime, err = time.Parse(layout, startUpdated)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid startUpdated date format")
		}
		endUpdatedTime, err = time.Parse(layout, endUpdated)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid endUpdated date format")
		}
	}

	servers, err := s.repository.GetServersByOptionalDateRange(startCreatedTime, endCreatedTime, startUpdatedTime, endUpdatedTime, field, order)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error fetching servers: %v", err))
	}

	f, err2 := util.CreateExcelFile(servers)

	if err2 != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error creating Excel file: %v", err2))
	}

	// Save Excel file to disk
	filePath := "export.xlsx"
	if err := f.SaveAs(filePath); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save Excel file")
	}

	// Serve the file
	return c.Attachment(filePath, "export.xlsx")
}

// GetServerUptime calculates the uptime for a server for the entire specified day.
func (s *Service) GetServerUptime(c echo.Context, serverID string, date string) (time.Duration, error) {
	layout := "2006-01-02"
	day, err := time.Parse(layout, date)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "Invalid date format: "+err.Error())
	}

	uptime, err := s.elastic.CalculateServerUptime(serverID, day)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, "Error calculating uptime: "+err.Error())
	}

	return uptime, nil
}

// GetServerReport sends a report of server statuses within a specified date range to the client.
func (s *Service) GetServerReport(c echo.Context, mail, start, end string) error {
	layout := "2006-01-02"
	location, err := time.LoadLocation("Asia/Bangkok") // Load the GMT+7 timezone

	if _, err := time.ParseInLocation(layout, start, location); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid start date format")
	}

	if _, err := time.ParseInLocation(layout, end, location); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid end date format")
	}

	mailArr := []string{mail}

	// Create a gRPC client
	var addr string = "127.0.0.1:50052" // Address of the gRPC server
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to dial server: %v", err)
	}

	defer conn.Close()

	client := pb.NewReportServiceClient(conn)

	err = doSendReport(client, mailArr, start, end)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error sending report: "+err.Error())
	}

	return c.String(http.StatusOK, "Report sent successfully")
}
