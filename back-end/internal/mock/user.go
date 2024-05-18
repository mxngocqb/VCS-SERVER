package mock

import (
	"github.com/labstack/echo/v4"
	model "github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (mockService *MockUserService) View(ctx echo.Context, id string) (*model.User, error) {
	args := mockService.Called(ctx, id)
    user, ok := args.Get(0).(*model.User)
    if !ok && user == nil {
        return nil, args.Error(1)
    }
    return user, args.Error(1)
}

func (mockService *MockUserService) Create(ctx echo.Context, u *model.User) (*model.User, error) {
	args := mockService.Called(ctx, u)
	user, ok := args.Get(0).(*model.User)
	if !ok && user == nil {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (mockService *MockUserService) Update(ctx echo.Context, id string, u *model.User) (*model.User, error) {
	args := mockService.Called(ctx, id, u)
	return args.Get(0).(*model.User), args.Error(1)
}

func (mockService *MockUserService) Delete(ctx echo.Context, id string) error {
	args := mockService.Called(ctx, id)
	return args.Error(0)
}
