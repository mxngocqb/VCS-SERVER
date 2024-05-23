package transport

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/handler/server"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/util"
)

type HTTP struct {
	service server.IServerService
}

// CreateRequest represents the request body for creating a server.
// @Summary Create server request
// @Description Represents the request body for creating a server.
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

// View gets a list of servers based on the provided filters and pagination.
// @Summary Get servers
// @Description Returns a list of servers based on the provided filters and pagination.
// @Tags Server
// @Accept json
// @Produce json
// @Param limit query int false "Limit of servers returned" default(10)
// @Param offset query int false "Ofset in server list" default(0)
// @Param status query string false "Filter by status"
// @Param field query string false "Field to sort by"
// @Param order query string false "Order by" Enums(asc, desc)
// @Success 200 {array} ServerResponse
// @Failure 400 {object} echo.HTTPError "Bad request - Invalid parameters for limit or offset or status or field or order"
// @Failure 404 {object} echo.HTTPError "Failed to fetch servers: No servers found based on the filters provided or server does not exist"
// @Failure 500 {object} echo.HTTPError "Internal server error - Failed to fetch servers"
// @Router /servers [get]
// @Security Bearer
func (h HTTP) View(c echo.Context) error {
	fmt.Println("View call")
	r := new(ViewRequest)

	// Retrieve servers from the database
	if err := c.Bind(r); err != nil {
		return err
	}

	servers, numberOfServers, err := h.service.View(c, r.Limit, r.Offset, r.Status, r.Field, r.Order)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Failed to fetch servers: "+err.Error())
	}

	response := ServerResponse{
		Total: numberOfServers,
		Data:  servers,
	}

	return c.JSON(http.StatusOK, response)
}

// Export exports filtered server data to an Excel file.
// @Summary Export servers to Excel
// @Description Exports server data to an Excel file based on the provided filters.
// @Tags Server
// @Produce application/octet-stream
// @Param limit query int false "Limit of servers returned" default(10)
// @Param offset query int false "Ofset in server list" default(0)
// @Param status query string false "Filter by status"
// @Param field query string false "Field to sort by"
// @Param order query string false "Order by" Enums(asc, desc)
// @Success 200 {file} file "Excel file containing server data"
// @Failure 400 {object} echo.HTTPError "Bad request - Invalid parameters for limit or offset or status"
// @Failure 404 {object} echo.HTTPError "Bad request - No servers found based on the filters provided or server does not exist"
// @Failure 403 {object} echo.HTTPError "Forbidden - User does not have permission to export servers"
// @Failure 409 {object} echo.HTTPError "Conflict - Failed to generate or send file"
// @Failure 500 {object} echo.HTTPError "Internal server error - Failed to generate or send file"
// @Router /servers/export [get]
// @Security Bearer
func (h HTTP) Export(c echo.Context) error {
	// Optional query parameters
	r := new(ViewRequest)

	// Retrieve servers from the database
	if err := c.Bind(r); err != nil {
		return err
	}

	err := h.service.GetServersFiltered(c, r.Limit, r.Offset, r.Status, r.Field, r.Order)
	if err != nil {
		return err
	}

	return nil
}

// Creeate creates a new server in the database.
// @Summary Create server
// @Description Creates a new server in the database based on the provided data.
// @Tags Server
// @Accept json
// @Produce json
// @Param server body CreateRequest true "Server data"
// @Success 201 {object} model.Server
// @Failure 400 {object} echo.HTTPError "Bad request - Invalid server data"
// @Failure 500 {object} echo.HTTPError "Internal server error - Failed to create serve"
// @Failure 403 {object} echo.HTTPError "Forbidden - User does not have permission to create server"
// @Router /servers [post]
// @Security Bearer
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
		return err
	}

	return c.JSON(http.StatusCreated, createdServer)
}

// Update updates a server in the database.
// @Summary Update server by ID
// @Description Updates a server in the database based on the provided ID.
// @Tags Server
// @Accept json
// @Produce json
// @Param id path int true "Server ID"
// @Param server body UpdateRequest true "Server update data"
// @Success 200 {object} model.Server
// @Failure 400 {object} echo.HTTPError "Bad request - Invalid update data"
// @Failure 404 {object} echo.HTTPError "Not found - Server not found"
// @Failure 500 {object} echo.HTTPError "Internal server error - Failed to update server"
// @Failure 403 {object} echo.HTTPError "Forbidden - User does not have permission to update server"
// @Router /servers/{id} [put]
// @Security Bearer
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
		return err
	}

	return c.JSON(http.StatusOK, u)
}

// Delete removes a server from the database.
// @Summary Delete server by ID
// @Description Deletes a server from the database based on the provided ID.
// @Tags Server
// @Accept json
// @Produce json
// @Param id path int true "Server ID"
// @Success 204
// @Failure 404 {object} echo.HTTPError "Not found - Not found server with the provided ID"
// @Failure 500 {object} echo.HTTPError "Internal server error - Failed to Delete server"
// @Failure 403 {object} echo.HTTPError "Forbidden - User does not have permission to delete server"
// @Router /servers/{id} [delete]
// @Security Bearer
func (h HTTP) Delete(c echo.Context) error {
	fmt.Println("Delete call")
	id := c.Param("id")
	err := h.service.Delete(c, id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// CreateMany creates multiple servers from an uploaded Excel file.
// @Summary Import servers from Excel
// @Description Imports server data from an Excel file and creates multiple servers.
// @Tags Server
// @Accept multipart/form-data
// @Produce json
// @Param listserver formData file true "Excel file containing server data"
// @Success 200 {object} ImportServerResponse
// @Failure 400 {object} echo.HTTPError "Bad request - Failed to read or open file"
// @Failure 500 {object} echo.HTTPError "Internal server error - Failed to parse Excel or create servers"
// @Failure 403 {object} echo.HTTPError "Forbidden - User does not have permission to import servers"
// @Router /servers/import [post]
// @Security Bearer
func (h HTTP) CreateMany(c echo.Context) error {
	file, err := c.FormFile("listserver")
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
		return err
	}

	response := ImportServerResponse{
		Message: "Import servers successfully",
		Total_success: len(createdServers),
		Lists_success: successLines,
		Total_fail: len(failedLines),
		Lists_fail: failedLines,
	}

	return c.JSON(http.StatusOK, response)
}


// GetServerUpTime get the uptime of a specified server.
// @Summary Get server uptime based on date and server ID
// @Description Returns the total hours of uptime for a specific server on a given date.
// @Tags Server
// @Accept json
// @Produce json
// @Param id path int true "Server ID" description("The ID of the server to get uptime for")
// @Param date query string true "Date" description("Formatted as YYYY-MM-DD, the date to get uptime for")
// @Success 200 {number} float64 "Total hours of uptime for the server on the specified date"
// @Failure 400 {object} echo.HTTPError "Invalid server ID or date format"
// @Failure 500 {object} echo.HTTPError "Internal server error occurred while fetching uptime"
// @Router /servers/{id}/uptime [get]
// @Security Bearer
func (h HTTP) GetServerUpTime(c echo.Context) error {
	serverID := c.Param("id")
	date := c.QueryParam("date")

	uptime, err := h.service.GetServerUptime(c, serverID, date)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, uptime.Hours())
}


// GetServersReport generates a report of server statuses within a specified date range.
// @Summary Send daily server report to administator email
// @Description Send a report of daily server statuses from the specified date range to the provided email.
// @Tags Server
// @Accept json
// @Produce json
// @Param start query string true "From Date" description("Formatted as YYYY-MM-DD, the start date of the report range")
// @Param end query string true "To Date" description("Formatted as YYYY-MM-DD, the end date of the report range")
// @Param mail query string true "Administrator Email" description("Administrator email to send the report to")
// @Success 200 {string} string "Sent report to administrator email successfully"
// @Failure 400 {object} echo.HTTPError "Invalid administrator email or date range"
// @Failure 500 {object} echo.HTTPError "Internal server error occurred while generating report"
// @Router /servers/report [get]
// @Security Bearer
func (h HTTP) GetServersReport(c echo.Context) error {

	start := c.QueryParam("start")
	end := c.QueryParam("end")
	mail := c.QueryParam("mail")

	err := h.service.GetServerReport(c, mail, start, end)
	if err != nil {
		return err
	}

	return nil
}
