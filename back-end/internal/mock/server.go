package mock

import (
	"time"

	"github.com/labstack/echo/v4"
	model "github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockServerService struct {
	mock.Mock
}

func (mockService *MockServerService) View(c echo.Context, limit, offset int, status, field, order string) ([]model.Server, int64, error) {
	args := mockService.Called(c, limit, offset, status, field, order)
	return args.Get(0).([]model.Server),1, args.Error(1)
}

func (mockService *MockServerService) Create(c echo.Context, server *model.Server) (*model.Server, error) {
	args := mockService.Called(c, server)
	return args.Get(0).(*model.Server), args.Error(1)
}

func (mockService *MockServerService) Update(c echo.Context, id string, server *model.Server) (*model.Server, error) {
	args := mockService.Called(c, id, server)
	return args.Get(0).(*model.Server), args.Error(1)
}

func (mockService *MockServerService) Delete(c echo.Context, id string) error {
	args := mockService.Called(c, id)
	return args.Error(0)
}

func (mockService *MockServerService) CreateMany(c echo.Context, servers []model.Server) ([]model.Server, []int, []int, error) {
	args := mockService.Called(c, servers)
	return args.Get(0).([]model.Server), args.Get(1).([]int), args.Get(2).([]int), args.Error(3)
}

func (mockService *MockServerService) GetServersFiltered(c echo.Context, startCreated, endCreated, startUpdated, endUpdated, field, order string) error {
	args := mockService.Called(c, startCreated, endCreated, startUpdated, endUpdated, field, order)
	return args.Error(0)
}

func (mockService *MockServerService) GetServerUptime(c echo.Context, serverID string, date string) (time.Duration, error) {
	args := mockService.Called(c, serverID, date)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (mockService *MockServerService) GetServerReport(c echo.Context, mail, start, end string) error {
	args := mockService.Called(c, mail, start, end)
	return args.Error(0)
}