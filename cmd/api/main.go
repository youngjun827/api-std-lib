package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/youngjun827/api-std-lib/internal/database"
	"github.com/youngjun827/api-std-lib/internal/database/models"
)

type application struct {
	logger *slog.Logger
	users  models.UserInterface
}

func main() {
	runtime.GOMAXPROCS(1)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ AddSource: true,}))

	db, err := database.InitDB()
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
			logger.Error("Could not start server: %v", err)
			os.Exit(1)
		}
	}()

	<-done
	logger.Info("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server Shutdown Failed: %v", err)
		os.Exit(1)
	}
	logger.Info("Server Exited Properly")
}
