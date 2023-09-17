package models

import (
	"database/sql"
	"errors"
	"log/slog"
)

type UserInterface interface {
	CreateUserQuery(user User) (int, error)
	GetUserByIDQuery(id int) (User, error)
	ListUsersQuery() ([]User, error)
	UpdateUserQuery(id int, user User) error
	DeleteUserQuery(id int) error
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

func (m *UserModel) CreateUserQuery(user User) (int, error) {
	sqlStatement := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := m.DB.QueryRow(sqlStatement, user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return 0, ErrNoModels
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) GetUserByIDQuery(id int) (User, error) {
	sqlStatement := `SELECT id, name, email, password FROM users WHERE id=$1`
	row := m.DB.QueryRow(sqlStatement, id)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoModels
		} else {
			return User{}, err
		}
	}
	return user, nil
}

func (m *UserModel) ListUsersQuery() ([]User, error) {
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

func (m *UserModel) UpdateUserQuery(id int, user User) error {
	sqlStatement := `UPDATE users SET name=$1, email=$2, password=$3 WHERE id=$4`
	_, err := m.DB.Exec(sqlStatement, user.Name, user.Email, user.Password, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return ErrNoModels
		} else {
			return err
		}
	}
	return nil
}

func (m *UserModel) DeleteUserQuery(id int) error {
	sqlStatement := `DELETE FROM users WHERE id=$1`
	_, err := m.DB.Exec(sqlStatement, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return ErrNoModels
		} else {
			return err
		}
	}
	return nil
}
