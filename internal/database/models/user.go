package models

import (
	"database/sql"
	"log/slog"
)

type UserInterface interface {
	CreateUser(user User) (int, error)
	GetUserByID(id int) (User, error)
	ListUsers() ([]User, error)
	UpdateUser(id int, user User) (error)
	DeleteUser(id int) error
}
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) CreateUser(user User) (int, error) {
	sqlStatement := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := m.DB.QueryRow(sqlStatement, user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		return 0, err
	}
	return id, nil
}

func (m *UserModel) GetUserByID(id int) (User, error) {
	sqlStatement := `SELECT id, name, email, password FROM users WHERE id=$1`
	row := m.DB.QueryRow(sqlStatement, id)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		slog.Error("Failed to get user by ID", "error", err)
		return User{}, err
	}
	return user, nil
}

func (m *UserModel) ListUsers() ([]User, error) {
	sqlStatement := `SELECT id, name, email, password FROM users`
	rows, err := m.DB.Query(sqlStatement)
	if err != nil {
		slog.Error("Failed to list users", "error", err)
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			slog.Error("Failed to scan user row", "error", err)
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		slog.Error("Failed to iterate over user rows", "error", err)
		return nil, err
	}
	return users, nil
}

func (m *UserModel) UpdateUser(id int, user User) (error) {
	sqlStatement := `UPDATE users SET name=$1, email=$2, password=$3 WHERE id=$4`
	_, err := m.DB.Exec(sqlStatement, user.Name, user.Email, user.Password, id)
	if err != nil {
		slog.Error("Failed to update user", "error", err)
		return err
	}
	return nil
}

func (m *UserModel) DeleteUser(id int) error {
	sqlStatement := `DELETE FROM users WHERE id=$1`
	_, err := m.DB.Exec(sqlStatement, id)
	if err != nil {
		slog.Error("Failed to delete user", "error", err)
		return err
	}
	return nil
}
