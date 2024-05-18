package mock

import (
	"errors"

	model "github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (mockService *MockAuthService) Authenticate(username, password string) (*model.User, error) {
	args := mockService.Called(username, password)
	if args.Get(0) == nil {
		return nil, errors.New("authentication failed")
	}
	return args.Get(0).(*model.User), args.Error(1)
}
