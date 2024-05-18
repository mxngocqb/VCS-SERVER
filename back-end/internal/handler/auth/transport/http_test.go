package transport

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	mockService "github.com/mxngocqb/VCS-SERVER/back-end/internal/mock"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	// Create a new instance of the repository mock
	mockSvc := new(mockService.MockAuthService)
	h := HTTP{service: mockSvc}

	// Create a new Echo instance
	e := echo.New()

	// Define a test case
	t.Run("successful authentication", func(t *testing.T) {
		// Set up expectations for the mock service
		mockSvc.On("Authenticate", "testuser", "testpassword").Return(&model.User{Username: "testuser", Password: "testpassword"}, nil)

		// Create a request to mimic login
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username": "testuser", "password": "testpassword"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the login handler
		err := h.Login(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "token")
		mockSvc.AssertExpectations(t)
	})

	t.Run("authentication failure", func(t *testing.T) {
		// Set up expectations for the mock service
		mockSvc.On("Authenticate", "testuser", "wrongpassword").Return(nil, errors.New("authentication failed"))

		// Create a request to mimic login
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username": "testuser", "password": "wrongpassword"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the login handler
		err := h.Login(c)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockSvc.AssertExpectations(t)
	})
}
