package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/handler"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/repository"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/cache"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type IService interface {
	View(c echo.Context, perPage int, offset int, status, field, order string) ([]model.Server, error)
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
	elastic    *service.ElasticService
	cache      cache.ServerCache
}

func NewServerService(repository *repository.ServerRepository, rbac *handler.RbacService, elastic *service.ElasticService, sc cache.ServerCache) *Service {
	return &Service{
		repository: repository,
		rbac:       rbac,
		elastic:    elastic,
		cache:      sc,
	}
}

// View retrieves servers from the database with optional pagination and status filtering.
func (s *Service) View(c echo.Context, perPage int, offset int, status, field, order string) ([]model.Server, error) {
	// ctx := c.Request().Context()
	key := s.cache.ConstructCacheKey(perPage, offset, status, field, order)

	// Try to get data from Redis first
	values := s.cache.GetMultiRequest(key)

	if values != nil {
		// Data found in cache
		return values, nil
	}

	// Data not found in cache, fetch from database
	servers, err := s.repository.GetServersFiltered(perPage, offset, status, field, order)
	if err != nil {
		return nil, err
	}

	s.cache.SetMultiRequest(key, servers) // Adjust expiration as needed
	return servers, nil
}

// Create creates a new server.
func (s *Service) Create(c echo.Context, server *model.Server) (*model.Server, error) {

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
	
	// Cache the server
	s.cache.Set(strconv.Itoa(int(server.ID)), server) 

	return server, nil
}

// CreateMany creates multiple servers and returns detailed results.
func (s *Service) CreateMany(c echo.Context, servers []model.Server) ([]model.Server, []int, []int, error) {
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

		// After successfully creating the server, log the status change
		err = s.elastic.LogStatusChange(server, server.Status)
		if err != nil {
			// Handle logging error
			fmt.Println("Error logging status change for server ID in Elasticsearch", server.ID)
		}
	}

	return createdServers, successLines, failedLines, nil
}

// Update updates a server.
func (s *Service) Update(c echo.Context, id string, server *model.Server) (*model.Server, error) {

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
	s.cache.Set(strconv.Itoa(int(server.ID)), server) 

	return updatedServer, nil
}

// Delete deletes a server.
func (s *Service) Delete(c echo.Context, id string) error {
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

	f := excelize.NewFile()
	index, err := f.NewSheet("Servers")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create Excel sheet")
	}
	f.SetActiveSheet(index)

	// Create header
	headers := []string{"ID", "Name", "Status", "IP", "Created At", "Updated At"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Servers", cell, header)
	}

	// Fill data
	for i, server := range servers {
		row := i + 2 // Starting from the second row
		f.SetCellValue("Servers", "A"+strconv.Itoa(row), server.ID)
		f.SetCellValue("Servers", "B"+strconv.Itoa(row), server.Name)
		f.SetCellValue("Servers", "C"+strconv.Itoa(row), server.Status)
		f.SetCellValue("Servers", "D"+strconv.Itoa(row), server.IP)
		f.SetCellValue("Servers", "E"+strconv.Itoa(row), server.CreatedAt.Format(time.RFC3339))
		f.SetCellValue("Servers", "F"+strconv.Itoa(row), server.UpdatedAt.Format(time.RFC3339))
	}

	// Save Excel file to disk
	filePath := "servers.xlsx"
	if err := f.SaveAs(filePath); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save Excel file")
	}

	// Serve the file
	return c.Attachment(filePath, "servers.xlsx")
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

	startTime, err := time.ParseInLocation(layout, start, location)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid start date format")
	}

	endTime, err := time.ParseInLocation(layout, end, location)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid end date format")
	}

	mailArr := []string{mail}

	err = service.SendReport(mailArr, startTime, endTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error sending report: "+err.Error())
	}

	return c.String(http.StatusOK, "Report sent successfully")
}
