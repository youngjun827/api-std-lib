package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	errEnv := loadEnvVariables()
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

func loadEnvVariables() error {
	// Open the .env file.
	file, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file line by line.
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Parse and set the environment variables.
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}

	return nil
}
