package db

import (
	"errors"
	"sync"

	"github.com/youngjun827/api-std-lib/api/models"
)

type UserRepositoryMock struct {
	MockCreateUser  func(user models.User) (int, error)
	MockGetUserByID func(id int) (models.User, error)
	MockListUsers   func() ([]models.User, error)
	MockUpdateUser  func(id int, user models.User) error
	MockDeleteUser  func(id int) error
	data            map[int]models.User
	mutex           *sync.Mutex
	lastID          int
}

func NewUserRepositoryMock() *UserRepositoryMock {
	return &UserRepositoryMock{
		data:  make(map[int]models.User),
		mutex: &sync.Mutex{},
		lastID: 0,
	}
}

func (m *UserRepositoryMock) CreateUser(user models.User) (int, error) {
	if m.MockCreateUser != nil {
		return m.MockCreateUser(user)
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.lastID++
	user.ID = m.lastID
	m.data[user.ID] = user
	return user.ID, nil
}

func (m *UserRepositoryMock) GetUserByID(id int) (models.User, error) {
	if m.MockGetUserByID != nil {
		return m.MockGetUserByID(id)
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	user, ok := m.data[id]
	if !ok {
		return models.User{}, errors.New("User not found")
	}
	return user, nil
}

func (m *UserRepositoryMock) ListUsers() ([]models.User, error) {
	if m.MockListUsers != nil {
		return m.MockListUsers()
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var users []models.User
	for _, user := range m.data {
		users = append(users, user)
	}
	return users, nil
}

func (m *UserRepositoryMock) UpdateUser(id int, user models.User) error {
	if m.MockUpdateUser != nil {
		return m.MockUpdateUser(id, user)
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, ok := m.data[id]
	if !ok {
		return errors.New("User not found")
	}
	user.ID = id
	m.data[id] = user
	return nil
}

func (m *UserRepositoryMock) DeleteUser(id int) error {
	if m.MockDeleteUser != nil {
		return m.MockDeleteUser(id)
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, ok := m.data[id]
	if !ok {
		return errors.New("User not found")
	}
	delete(m.data, id)
	return nil
}
