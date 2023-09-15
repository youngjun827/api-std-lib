package database

import (
	"database/sql"

	"log/slog"

	"github.com/youngjun827/api-std-lib/internal/database/models"
)

type UserRepository interface {
	CreateUser(user models.User) (int, error)
	GetUserByID(id int) (models.User, error)
	ListUsers() ([]models.User, error)
	UpdateUser(id int, user models.User) error
	DeleteUser(id int) error
}

type UserRepositorySQL struct {
	DB *sql.DB
}

func NewUserRepositorySQL(db *sql.DB) *UserRepositorySQL {
	return &UserRepositorySQL{DB: db}
}

func (ur *UserRepositorySQL) CreateUser(user models.User) (int, error) {
	sqlStatement := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := ur.DB.QueryRow(sqlStatement, user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		return 0, err
	}
	return id, nil
}

func (ur *UserRepositorySQL) GetUserByID(id int) (models.User, error) {
	sqlStatement := `SELECT id, name, email, password FROM users WHERE id=$1`
	row := ur.DB.QueryRow(sqlStatement, id)
	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		slog.Error("Failed to get user by ID", "error", err)
		return models.User{}, err
	}
	return user, nil
}

func (ur *UserRepositorySQL) ListUsers() ([]models.User, error) {
	sqlStatement := `SELECT id, name, email, password FROM users`
	rows, err := ur.DB.Query(sqlStatement)
	if err != nil {
		slog.Error("Failed to list users", "error", err)
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var user models.User
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

func (ur *UserRepositorySQL) UpdateUser(id int, user models.User) error {
	sqlStatement := `UPDATE users SET name=$1, email=$2, password=$3 WHERE id=$4`
	_, err := ur.DB.Exec(sqlStatement, user.Name, user.Email, user.Password, id)
	if err != nil {
		slog.Error("Failed to update user", "error", err)
		return err
	}
	return nil
}

func (ur *UserRepositorySQL) DeleteUser(id int) error {
	sqlStatement := `DELETE FROM users WHERE id=$1`
	_, err := ur.DB.Exec(sqlStatement, id)
	if err != nil {
		slog.Error("Failed to delete user", "error", err)
		return err
	}
	return nil
}
