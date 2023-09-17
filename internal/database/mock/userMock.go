package mock

import (
	"errors"

	"github.com/youngjun827/api-std-lib/internal/database/models"
)

var MockUser = models.User{
	ID:       1,
	Name:     "Test",
	Email:    "test@example.com",
	Password: "Jooa005500!",
}

type MockUserModel struct {
	NoMatchErrorMode string
	ServerErrorMode string
}

func NewMockUserModel(errorMode string) *MockUserModel {
	if errorMode == "NoMatch" {
		return &MockUserModel{
			NoMatchErrorMode: errorMode,
		}
	} 
	if errorMode == "ServerError" {
		return &MockUserModel{
			ServerErrorMode: errorMode,
		}
	} 	
	return &MockUserModel{}
}

func (m *MockUserModel) CreateUserQuery(user models.User) (int, error) {
	if m.NoMatchErrorMode == "NoMatch" {
		if user.Email == "exists@example.com" {
			return 0, models.ErrNoModels
		}
	}
	if m.ServerErrorMode == "ServerError" {
		return 0, errors.New("unexpected error")
	}
	return MockUser.ID, nil
}

func (m *MockUserModel) GetUserByIDQuery(id int) (models.User, error) {
	if m.NoMatchErrorMode == "NoMatch" {
		return models.User{}, models.ErrNoModels
	}
	if m.ServerErrorMode == "ServerError" {
		return models.User{}, errors.New("unexpected error")
	}
	return MockUser, nil
}

func (m *MockUserModel) ListUsersQuery() ([]models.User, error) {
	if m.ServerErrorMode == "ServerError" {
		return nil, errors.New("unexpected error")
	}
	if MockUser.ID == 0 {
		return nil, errors.New("no users found")
	}

	return []models.User{MockUser}, nil
}

func (m *MockUserModel) UpdateUserQuery(id int, user models.User) error {
	if m.NoMatchErrorMode == "NoMatch" {
		return models.ErrNoModels
	}
	if m.ServerErrorMode == "ServerError" {
		return errors.New("unexpected error")
	}
	return nil
}

func (m *MockUserModel) DeleteUserQuery(id int) error {
	if m.NoMatchErrorMode == "NoMatch" {
		return models.ErrNoModels
	}
	if m.ServerErrorMode == "ServerError" {
		return errors.New("unexpected error")
	}
	return nil
}
