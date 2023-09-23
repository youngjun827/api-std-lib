package main

import (
	"bufio"
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/youngjun827/api-std-lib/internal/database/models"
	error_response "github.com/youngjun827/api-std-lib/internal/error"
	"github.com/youngjun827/api-std-lib/internal/validator"
)

type application struct {
	logger *slog.Logger
	users  models.UserInterface
	validator.Validator
	error_response.ErrorResponse
}

func main() {
	runtime.GOMAXPROCS(1)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	db, err := initDB()
	if err != nil {
		logger.Error("Database Connection Refused", "error", err)
	}

	app := &application{
		logger: logger,
		users:  &models.UserModel{DB: db},
	}

	logger.Info("Starting Server on PORT:8081")

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      app.SetupRoutes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Could not start server", "error", err)
			os.Exit(1)
		}
	}()

	<-done
	logger.Info("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server Shutdown Failed", "error", err)
		os.Exit(1)
	}
	logger.Info("Server Exited Properly")
}

func initDB() (*sql.DB, error) {
	errEnv := loadEnvVariables()
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

func loadEnvVariables() error {
	file, err := os.Open(".env")
	if err != nil {
		slog.Error("Failed to load the environment variable .env", "error", err)
		return err
	}
	defer file.Close()
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
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
