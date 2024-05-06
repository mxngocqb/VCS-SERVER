package transport

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/handler/server"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/util"
	"net/http"
)

type HTTP struct {
	service server.IService
}

// NewHTTP sets up server-related routes.
// @Summary Initialize server routes
// @Description Configures routes for server operations such as CRUD and import/export functionalities.
// @Tags Server
func NewHTTP(r *echo.Group, service *server.Service) {
	h := HTTP{service}

	sr := r.Group("/servers")
	sr.POST("", h.Create)
	sr.POST("/import", h.CreateMany)
	sr.GET("", h.View)
	sr.PUT("/:id", h.Update)
	sr.DELETE("/:id", h.Delete)
	sr.GET("/export", h.Export)
	sr.GET("/:id/uptime", h.GetServerUpTime)
	sr.GET("/report", h.GetServersReport)
}

// View retrieves a list of servers with optional pagination and filters.
// @Summary List servers
// @Description Retrieves servers based on provided pagination and optional filters.
// @Tags Server
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of servers returned" default(10)
// @Param offset query int false "Offset in server list" default(0)
// @Param status query string false "Filter by status"
// @Param field query string false "Field to sort by"
// @Param order query string false "Order of sort" Enums(asc, desc)
// @Success 200 {array} model.Server
// @Failure 400 {object} echo.HTTPError "Invalid parameters for limit or offset"
// @Failure 500 {object} echo.HTTPError "Failed to fetch servers due to server error"
// @Router /servers [get]
func (h HTTP) View(c echo.Context) error {
	fmt.Println("View call")
	r := new(ViewRequest)

	// Retrieve servers from the database
	if err := c.Bind(r); err != nil {
		return err
	}

	fmt.Println(r)

	servers, err := h.service.View(c, r.Limit, r.Offset, r.Status, r.Field, r.Order)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch servers: "+err.Error())
	}

	return c.JSON(http.StatusOK, servers)
}

// Create adds a new server to the database.
// @Summary Create a server
// @Description Adds a new server to the database.
// @Tags Server
// @Accept json
// @Produce json
// @Param server body CreateRequest true "Server data"
// @Success 201 {object} model.Server
// @Failure 400 {object} echo.HTTPError "Invalid server data provided"
// @Failure 500 {object} echo.HTTPError "Failed to create server due to server error"
// @Router /servers [post]
func (h HTTP) Create(c echo.Context) error {
	r := new(CreateRequest)
	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Create a new server model
	newServer := &model.Server{
		Name:   r.Name,
		Status: r.Status,
		IP:     r.IP,
	}

	// Create new server in the database
	createdServer, err := h.service.Create(c, newServer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, createdServer)
}

// Update modifies an existing server.
// @Summary Update a server
// @Description Updates server details.
// @Tags Server
// @Accept json
// @Produce json
// @Param id path int true "Server ID"
// @Param server body UpdateRequest true "Server update data"
// @Success 200 {object} model.Server
// @Failure 400 {object} echo.HTTPError "Bad request - Invalid update data"
// @Failure 404 {object} echo.HTTPError "Not found - Server not found"
// @Failure 500 {object} echo.HTTPError "Internal server error - Failed to update server"
// @Router /servers/{id} [put]
func (h HTTP) Update(c echo.Context) error {
	// Retrieve the request body
	r := new(UpdateRequest)
	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Retrieve the server ID from the URL
	id := c.Param("id")

	// Create a new server model
	updatedServer := &model.Server{
		Name:   r.Name,
		Status: r.Status,
		IP:     r.IP,
	}

	// Update server in the database
	u, err := h.service.Update(c, id, updatedServer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, u)
}

// Delete removes a server from the database.
// @Summary Delete a server
// @Description Removes a server based on ID.
// @Tags Server
// @Accept json
// @Produce json
// @Param id path int true "Server ID"
// @Success 204
// @Failure 404 {object} echo.HTTPError "Not found - Server not found"
// @Failure 500 {object} echo.HTTPError "Internal server error - Failed to delete server"
// @Router /servers/{id} [delete]
func (h HTTP) Delete(c echo.Context) error {
	fmt.Println("Delete call")
	id := c.Param("id")
	err := h.service.Delete(c, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// CreateMany handles the creation of multiple servers from an uploaded file.
// @Summary Bulk create servers
// @Description Creates multiple servers from an uploaded Excel file.
// @Tags Server
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel file with server data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} echo.HTTPError "Bad request - Invalid or corrupt file"
// @Failure 500 {object} echo.HTTPError "Internal server error - Failed to parse or save servers"
// @Router /servers/import [post]
func (h HTTP) CreateMany(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to read file: "+err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open file: "+err.Error())
	}
	defer src.Close()

	servers, err := util.ParseExcel(src)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse Excel: "+err.Error())
	}

	createdServers, successLines, failedLines, err := h.service.CreateMany(c, servers)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save servers: "+err.Error())
	}

	response := echo.Map{
		"message":       "Servers upload completed with detailed results.",
		"success_count": len(createdServers),
		"failure_count": len(failedLines),
		"success_lines": successLines,
		"failure_lines": failedLines,
	}

	return c.JSON(http.StatusOK, response)
}

// Export generates an Excel file of servers based on filters and sends it to the client.
// @Summary Export servers
// @Description Exports filtered server data to an Excel file.
// @Tags Server
// @Produce application/octet-stream
// @Param startCreated query string false "Filter by creation date start"
// @Param endCreated query string false "Filter by creation date end"
// @Param startUpdated query string false "Filter by update date start"
// @Param endUpdated query string false "Filter by update date end"
// @Param field query string false "Field to sort by"
// @Param order query string false "Order of sort"
// @Success 200 {file} file "Excel file"
// @Failure 400 {object} echo.HTTPError "Bad request - Invalid filter parameters"
// @Failure 500 {object} echo.HTTPError "Internal server error - Failed to generate or send file"
// @Router /servers/export [get]
func (h HTTP) Export(c echo.Context) error {
	// Optional query parameters
	startCreated := c.QueryParam("startCreated")
	endCreated := c.QueryParam("endCreated")
	startUpdated := c.QueryParam("startUpdated")
	endUpdated := c.QueryParam("endUpdated")
	field := c.QueryParam("field")
	order := c.QueryParam("order")

	err := h.service.GetServersFiltered(c, startCreated, endCreated, startUpdated, endUpdated, field, order)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

// GetServerUpTime retrieves the uptime of a specified server.
// @Summary Retrieve server uptime
// @Description Returns the uptime of a server based on a specific date provided in the query.
// @Tags Server
// @Accept json
// @Produce json
// @Param id path int true "Server ID" description("Unique identifier of the server")
// @Param date query string true "Date" description("The specific date to get the server uptime for, formatted as YYYY-MM-DD")
// @Success 200 {number} float64 "Hours of uptime"
// @Failure 400 {object} echo.HTTPError "Invalid date format or server ID"
// @Failure 500 {object} echo.HTTPError "Internal server error occurred while retrieving uptime"
// @Router /servers/{id}/uptime [get]
func (h HTTP) GetServerUpTime(c echo.Context) error {
	serverID := c.Param("id")
	date := c.QueryParam("date")

	uptime, err := h.service.GetServerUptime(c, serverID, date)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, uptime.Hours())
}

// GetServersReport generates a report of server statuses within a specified date range.
// @Summary Generate server status report
// @Description Retrieves a report of server statuses for a given date range and sends it to the specified email address.
// @Tags Server
// @Accept json
// @Produce json
// @Param mail body string true "Recipient Email" description("Email address to send the report to")
// @Param start query string true "Start Date" description("The start date of the report range, formatted as YYYY-MM-DD")
// @Param end query string true "End Date" description("The end date of the report range, formatted as YYYY-MM-DD")
// @Success 200 {string} string "Report sent successfully"
// @Failure 400 {object} echo.HTTPError "Invalid date format or email"
// @Failure 500 {object} echo.HTTPError "Error occurred while sending the report"
// @Router /servers/report [get]
func (h HTTP) GetServersReport(c echo.Context) error {
	r := new(GetServersReportRequest)
	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	start := c.QueryParam("start")
	end := c.QueryParam("end")

	err := h.service.GetServerReport(c, r.Mail, start, end)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}