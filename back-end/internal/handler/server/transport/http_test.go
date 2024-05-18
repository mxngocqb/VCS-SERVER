package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	mockService "github.com/mxngocqb/VCS-SERVER/back-end/internal/mock"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// TestHTTP_Create tests the Create method of the HTTP handler.
func TestHTTP_Create(t *testing.T) {
	// Setup the echo server
	e := echo.New()
	// Create a new server data
	server := model.Server{Name: "Create Server", Status: true, IP: "192.168.88.1"}
	body, _ := json.Marshal(server)
	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/servers", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Create a new mock service
	mockSvc := new(mockService.MockServerService)
	h := HTTP{service: mockSvc}
	// Mock the Create method
	mockSvc.On("Create", mock.Anything, mock.AnythingOfType("*model.Server")).Return(&server, nil)

	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		mockSvc.AssertExpectations(t)
	}
}

func TestHTTP_Update(t *testing.T) {
	// Setup the echo server
	e := echo.New()

	// Create a new server data
	mockService := new(mockService.MockServerService)
	handler := HTTP{service: mockService}

	// Create a new server data
	serverToUpdate := &model.Server{Name: "Updated Server", Status: true, IP: "192.168.88.2"}
	serverUpdated := &model.Server{Model: gorm.Model{ID: 1}, Name: "Updated Server", Status: true, IP: "192.168.88.2"}

	mockService.On("Update", mock.Anything, "1", serverToUpdate).Return(serverUpdated, nil)

	// Create a new HTTP request
	body, _ := json.Marshal(serverToUpdate)
	req := httptest.NewRequest(http.MethodPut, "/servers/1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/servers/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	assert.NoError(t, handler.Update(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

// TestHTTP_Delete tests the Delete method of the HTTP handler.
func TestHTTP_Delete(t *testing.T) {
	// Setup the echo server
	e := echo.New()
	mockService := new(mockService.MockServerService)
	handler := HTTP{service: mockService}

	// Mock the Delete method
	mockService.On("Delete", mock.Anything, "1").Return(nil)

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodDelete, "/servers/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/servers/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	assert.NoError(t, handler.Delete(c))            // Check if the handler returns no error
	assert.Equal(t, http.StatusNoContent, rec.Code) // Check if the response status code is 204
	mockService.AssertExpectations(t)
}

// TestHTTP_View tests the View method of the HTTP handler.
func TestHTTP_View(t *testing.T) {
	e := echo.New()
	e.Validator = &util.CustomValidator{Validator: validator.New()}
	e.Binder = &util.CustomBinder{Binder: &echo.DefaultBinder{}}
	handler := HTTP{service: new(mockService.MockServerService)}
	var nilError = (*echo.HTTPError)(nil)

	tests := []struct {
		name           string
		queryParams    string
		mockReturn     []model.Server
		mockError      error
		expectedStatus int
		expectedBody   string
		expected       *echo.HTTPError
	}{
		{
			name:        "Valid Request",
			queryParams: "?limit=10&offset=0",
			mockReturn: []model.Server{
				{Name: "Server1", Status: true, IP: "192.168.88.1"},
				{Name: "Server2", Status: false, IP: "192.168.88.2"},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Server1",
			expected:       nilError,
		},
		{
			name:        "Invalid Limit",
			queryParams: "?limit=-1&offset=0",
			expected:    echo.NewHTTPError(http.StatusBadRequest, "Key: 'ViewRequest.Limit' Error:Field validation for 'Limit' failed on the 'gte' tag"),
		},
		{
			name:           "Invalid Offset",
			queryParams:    "?limit=10&offset=-1&status=true&field=name&order=asc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Key: 'ViewRequest.Offset' Error:Field validation for 'Offset' failed on the 'gte' tag",
			expected:       echo.NewHTTPError(http.StatusBadRequest, "Key: 'ViewRequest.Offset' Error:Field validation for 'Offset' failed on the 'gte' tag"),
		},
		{
			name:        "Service Error",
			queryParams: "?limit=10&offset=0",
			mockError:   errors.New("internal error"),
			expected:    echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch servers: internal error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/servers"+tc.queryParams, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockSvc := new(mockService.MockServerService)
			mockSvc.On("View", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.mockReturn, tc.mockError)
			handler.service = mockSvc

			got := handler.View(c)
			if got != nil {
				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("expected %v, got %v", tc.expected, got)
				}
			}
		})
	}
}

// TestHTTP_CreateMany tests the CreateMany method of the HTTP handler.
func TestHTTP_CreateMany(t *testing.T) {
	// Setup the echo server
	e := echo.New()
	mockService := new(mockService.MockServerService)
	handler := HTTP{service: mockService}

	// Setup the multipart form data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Create a form file
	fw, err := w.CreateFormFile("file", "import_30.xlsx")
	assert.NoError(t, err)
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)
	f.SetCellValue(sheet, "A1", "Name")
    f.SetCellValue(sheet, "B1", "Status")
    f.SetCellValue(sheet, "C1", "IP")
    f.SetCellValue(sheet, "A2", "Server1")
    f.SetCellValue(sheet, "B2", "true")
    f.SetCellValue(sheet, "C2", "192.168.1.1")
	// Write fake content to the form file
	var excelBuffer bytes.Buffer
    err = f.Write(&excelBuffer)
    assert.NoError(t, err)
	// Write the form file to the form data
	_, err = fw.Write(excelBuffer.Bytes())
    assert.NoError(t, err)
    w.Close()
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/servers/import", &b)
	req.Header.Set("Content-Type",  w.FormDataContentType())
	
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Prepare the expected results from the service call
	expectedServers := []model.Server{{Name: "Server1", Status: true, IP: "192.168.1.1"}}
	expectedSuccess := []int{1}
	expectedFailures := []int{}

	// Setup expectations on the mock service
	mockService.On("CreateMany", mock.Anything, mock.AnythingOfType("[]model.Server")).Return(expectedServers, expectedSuccess, expectedFailures, nil)

	// Call the handler
	if assert.NoError(t, handler.CreateMany(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Here, further parsing of response JSON can be done to verify response content
	}

	// Assert that the expectations were met
	mockService.AssertExpectations(t)
}

func TestHTTP_Export(t *testing.T) {
    // Setup the echo server
    e := echo.New()
    mockService := new(mockService.MockServerService)
    handler := HTTP{service: mockService}

    // Define the query parameters
    queryParams := map[string]string{
        "startCreated": "2023-01-01",
        "endCreated":   "2023-12-31",
        "startUpdated": "2023-01-01",
        "endUpdated":   "2023-12-31",
        "field":        "name",
        "order":        "asc",
    }

    // Create a new HTTP request with the query parameters
    req := httptest.NewRequest(http.MethodGet, "/servers/export", nil)
    q := req.URL.Query()
    for key, value := range queryParams {
        q.Add(key, value)
    }
    req.URL.RawQuery = q.Encode()

    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    // Setup the mock service expectation
    mockService.On("GetServersFiltered", c, queryParams["startCreated"], queryParams["endCreated"], queryParams["startUpdated"], queryParams["endUpdated"], queryParams["field"], queryParams["order"]).Return(nil)

    // Call the handler
    if assert.NoError(t, handler.Export(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
    }

    // Assert that the expectations were met
    mockService.AssertExpectations(t)
}

func TestHTTP_GetServerUpTime(t *testing.T) {
    // Setup the echo server
    e := echo.New()
    mockService := new(mockService.MockServerService)
    handler := HTTP{service: mockService}

    // Define test parameters
    serverID := "1"
    date := "2024-05-17"

    // Prepare mock service expectations
    expectedUptime := 10 * time.Hour // Set an example duration
    mockService.On("GetServerUptime", mock.Anything, serverID, date).Return(expectedUptime, nil)

    // Create a new HTTP request
    req := httptest.NewRequest(http.MethodGet, "/servers/"+serverID+"/uptime?date="+date, nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(serverID)

    // Call the handler
    if assert.NoError(t, handler.GetServerUpTime(c)) {
        // Assert response status code
        assert.Equal(t, http.StatusOK, rec.Code)

        // Parse response body
        var uptime float64
        err := json.Unmarshal(rec.Body.Bytes(), &uptime)
        assert.NoError(t, err)

        // Assert uptime value
        assert.Equal(t, expectedUptime.Hours(), uptime)
    }

    // Assert that the expectations were met
    mockService.AssertExpectations(t)
}

func TestHTTP_GetServersReport(t *testing.T) {
    // Setup the echo server
    e := echo.New()
    mockService := new(mockService.MockServerService)
    handler := HTTP{service: mockService}

    // Define test parameters
    start := "2024-05-01"
    end := "2024-05-31"
    mail := "example@example.com"

    // Prepare mock service expectations
    mockService.On("GetServerReport", mock.Anything, mail, start, end).Return(nil)

    // Create a new HTTP request
    req := httptest.NewRequest(http.MethodGet, "/servers/report?start="+start+"&end="+end+"&mail="+mail, nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    // Call the handler
    if assert.NoError(t, handler.GetServersReport(c)) {
        // Assert response status code
        assert.Equal(t, http.StatusOK, rec.Code)
    }

    // Assert that the expectations were met
    mockService.AssertExpectations(t)
}
