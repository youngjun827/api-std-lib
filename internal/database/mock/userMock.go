package mock

import (
	"database/sql"
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
	ErrorMode bool
}

func NewMockUserModel(errorMode bool) *MockUserModel {
	return &MockUserModel{
		ErrorMode: errorMode,
	}
}

func (m *MockUserModel) CreateUserQuery(user models.User) (int, error) {
	if m.ErrorMode {
		if user.Email == "exists@example.com" {
			return 0, sql.ErrNoRows
		}
		return 0, errors.New("unexpected error")
	}
	return MockUser.ID, nil
}

func (m *MockUserModel) GetUserByIDQuery(id int) (models.User, error) {
	if m.ErrorMode {
		return models.User{}, errors.New("unexpected error")
	}
	if id == MockUser.ID {
		return MockUser, nil
	}
	return models.User{}, sql.ErrNoRows
}

func (m *MockUserModel) ListUsersQuery() ([]models.User, error) {
	if m.ErrorMode {
		return nil, errors.New("unexpected error")
	}
	if MockUser.ID == 0 {
		return nil, errors.New("no users found")
	}

	return []models.User{MockUser}, nil
}

func (m *MockUserModel) UpdateUserQuery(id int, user models.User) error {
	if m.ErrorMode {
		return errors.New("unexpected error")
	}
	if id != MockUser.ID {
		return sql.ErrNoRows
	}
	return nil
}

func (m *MockUserModel) DeleteUserQuery(id int) error {
	if m.ErrorMode {
		return errors.New("unexpected error")
	}
	if id != MockUser.ID {
		return sql.ErrNoRows
	}
	return nil
}
