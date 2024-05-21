package transport

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"github.com/labstack/echo/v4"
	mockService "github.com/mxngocqb/VCS-SERVER/back-end/internal/mock"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestHTTP_Create(t *testing.T) {
	// Setup the echo server
	e := echo.New()
	
	// Valid user data
	user := model.User{
		Model:    gorm.Model{},
		Username: "test",
		Password: "test",
		Roles:    []model.Role{},
		RoleIDs:  []uint{},
	}
	validBody, _ := json.Marshal(user)

	// Create a new mock service
	mockServices := new(mockService.MockUserService)
	h := HTTP{service: mockServices}

	t.Run("successful creation", func(t *testing.T) {
		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(validBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		// Mock the service
		mockServices.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(&user, nil)

		// Call the handler
		if assert.NoError(t, h.Create(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			mockServices.AssertExpectations(t)
		}
	})

	t.Run("invalid JSON input", func(t *testing.T) {
		invalidBody := []byte(`{"invalid json}`)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(invalidBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Create(c)
		if assert.Error(t, err) {
			he, ok := err.(*echo.HTTPError)
			if assert.True(t, ok) {
				assert.Equal(t, http.StatusBadRequest, he.Code)
			}
		}
	})
}


func TestHTTP_Delete(t *testing.T) {
    // Setup the echo server
    e := echo.New()

    t.Run("successful deletion", func(t *testing.T) {
        // Create a new HTTP request
        req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/users/:id")
        c.SetParamNames("id")
        c.SetParamValues("1")

        // Create a new mock service
        mockServices := new(mockService.MockUserService)
        handler := HTTP{service: mockServices}

        // Mock the Delete method
        mockServices.On("Delete", mock.Anything, "1").Return(nil)

        // Call the handler
        err := handler.Delete(c)
        if assert.NoError(t, err) {
            assert.Equal(t, http.StatusNoContent, rec.Code) // Check if the response status code is 204
            mockServices.AssertExpectations(t)
        }
    })

    t.Run("failed deletion", func(t *testing.T) {
		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodDelete, "/users/2", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("3")
	
		// Create a new mock service
		mockServices := new(mockService.MockUserService)
		handler := HTTP{service: mockServices}
	
		// Mock the Delete method to return an error
		mockErr := echo.NewHTTPError(http.StatusNotFound, "user not found")
		mockServices.On("Delete", mock.Anything, "3").Return(mockErr)
	
		// Call the handler
		err := handler.Delete(c)
		log.Println("Err: ",err)

		if assert.Error(t, err) {
			he, ok := err.(*echo.HTTPError)
			if assert.True(t, ok) {
				assert.Equal(t, http.StatusNotFound, he.Code) // Check if the response status code is 500
				assert.Contains(t, he.Message.(string), "not found")  // Check if the response body contains the error message
				mockServices.AssertExpectations(t)
			}
		}

	})
	
}


func TestHTTP_Update(t *testing.T) {
	// Setup the echo server
	e := echo.New()

	// Create a new server data
	mockService := new(mockService.MockUserService)
	handler := HTTP{service: mockService}

	// Create a new server data
	userToUpdate := &model.User{Username: "test", Password: "test"}
	userUpdated := &model.User{Model: gorm.Model{ID: 1}, Username: "test", Password: "test2"}

	mockService.On("Update", mock.Anything, "1", userToUpdate).Return(userUpdated, nil)

	// Create a new HTTP request
	body, _ := json.Marshal(userToUpdate)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	assert.NoError(t, handler.Update(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestViewHandler(t *testing.T) {
    // Create a new instance of our mock service
    mockService := new(mockService.MockUserService)

    // Create a new echo instance
    e := echo.New()
    // Setup the expected return values
    user := &model.User{Model: gorm.Model{ID: 1},Username: "Test User", Password: "password"}
    mockService.On("View", mock.Anything, "1").Return(user, nil)

    // Create a new HTTP request and recorder
    req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
    rec := httptest.NewRecorder()

    // Create a new echo context
    ctx := e.NewContext(req, rec)
    ctx.SetParamNames("id")
    ctx.SetParamValues("1")

    // Create the handler and call the method
    h := HTTP{service: mockService}
    if assert.NoError(t, h.View(ctx)) {
        assert.Equal(t, http.StatusOK, rec.Code)
        assert.JSONEq(t, `{
            "ID": 1,
            "username": "Test User",
            "password": "password",
            "CreatedAt": "0001-01-01T00:00:00Z",
            "UpdatedAt": "0001-01-01T00:00:00Z",
            "DeletedAt": null,
            "role_ids": null,
            "roles": null
        }`, rec.Body.String())
    }

    // Ensure that the expectations were met
    mockService.AssertExpectations(t)
}

func TestViewHandler_Error(t *testing.T) {
    // Create a new instance of our mock service
    mockService := new(mockService.MockUserService)


    // Create a new echo instance
    e := echo.New()

    // Setup the expected return values
    mockService.On("View", mock.Anything, "1").Return(nil, assert.AnError)

    // Create a new HTTP request and recorder
    req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
    rec := httptest.NewRecorder()

    // Create a new echo context
    ctx := e.NewContext(req, rec)
    ctx.SetParamNames("id")
    ctx.SetParamValues("1")

    // Create the handler and call the method
    h := HTTP{service: mockService}
    err := h.View(ctx)

    if assert.Error(t, err) {
        he, ok := err.(*echo.HTTPError)
        if assert.True(t, ok) {
            assert.Equal(t, http.StatusInternalServerError, he.Code)
            assert.Equal(t, assert.AnError.Error(), he.Message)
        }
    }

    // Ensure that the expectations were met
    mockService.AssertExpectations(t)
}
