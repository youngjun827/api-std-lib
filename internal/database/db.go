package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/youngjun827/api-std-lib/internal/middleware"
)

var DB *sql.DB

func InitDB() {
	errEnv := middleware.LoadEnvVariables()
	if errEnv != nil {
		slog.Error("Error loading .env file")
	}

	connStr := os.Getenv("DB_SOURCE")

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("Failed to connect to database")
	}

	// Set connection pool parameters
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(time.Minute)

	err = DB.Ping()
	if err != nil {
		slog.Error("Failed to ping database")
	}

	fmt.Println("Successfully connected to the database.")
}
