package mocks

import (
	"github.com/masudcsesust04/golang-jwt-auth/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock type for the UserDBInterface type
type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetAllUsers() ([]*models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockDB) GetUserByID(id int64) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDB) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDB) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDB) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDB) DeleteUser(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDB) CreateRefreshToken(refreshToken *models.RefreshToken) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

func (m *MockDB) GetRefreshToken(userID int64) (*models.RefreshToken, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockDB) DeleteRefreshToken(userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}