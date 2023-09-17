package database

import (
	"database/sql"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/youngjun827/api-std-lib/internal/middleware"
)

func InitDB() (*sql.DB, error) {
	errEnv := middleware.LoadEnvVariables()
	if errEnv != nil {
		slog.Error("Error loading .env file", "error", errEnv)
		return nil, errEnv
	}

	connStr := os.Getenv("DB_SOURCE")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		slog.Error("Failed to ping the database", "error", err)
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute)

	return db, nil
}
