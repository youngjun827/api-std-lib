package mock

import (
	"database/sql"
	"errors"

	"github.com/youngjun827/api-std-lib/internal/database/models"
)

const (
	NoMatch     = "NoMatch"
	ServerError = "ServerError"
)

var MockUser = models.User{
	ID:       1,
	Name:     "Test",
	Email:    "test@example.com",
	Password: "Jooa005500!",
}

type MockUserModel struct {
	ErrorMode string
}

func NewMockUserModel(errorMode string) *MockUserModel {
	return &MockUserModel{
		ErrorMode: errorMode,
	}
}

func (m *MockUserModel) createUserError() error {
	if m.ErrorMode == NoMatch {
		return models.ErrNoModels
	}
	if m.ErrorMode == ServerError {
		return errors.New("unexpected error")
	}
	return nil
}

func (m *MockUserModel) CreateUserQuery(user models.User) (int, error) {
	if err := m.createUserError(); err != nil {
		return 0, err
	}
	return MockUser.ID, nil
}

func (m *MockUserModel) GetUserByIDQuery(id int) (models.User, error) {
	if err := m.createUserError(); err != nil {
		return models.User{}, err
	}
	return MockUser, nil
}

func (m *MockUserModel) ListUsersQuery() ([]models.User, error) {
	if err := m.createUserError(); err != nil {
		return nil, err
	}
	return []models.User{MockUser}, nil
}

func (m *MockUserModel) UpdateUserQuery(id int, user models.User) error {
	if err := m.createUserError(); err != nil {
		return err
	}
	if id != MockUser.ID {
		return sql.ErrNoRows
	}
	return nil
}

func (m *MockUserModel) DeleteUserQuery(id int) error {
	if err := m.createUserError(); err != nil {
		return err
	}
	if id != MockUser.ID {
		return sql.ErrNoRows
	}
	return nil
}
