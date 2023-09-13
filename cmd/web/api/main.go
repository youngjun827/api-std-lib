package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/youngjun827/api-std-lib/cmd/web/api/routes"
	"github.com/youngjun827/api-std-lib/internal/database"
)

func main() {
	runtime.GOMAXPROCS(1)

	database.InitDB()

	userRepository := database.NewUserRepositorySQL(database.DB)

	mux := routes.SetupRoutes(userRepository)

	fmt.Println("Server is running on port 8081...")

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	slog.NewLogLogger(slog.Default().Handler(), slog.LevelInfo)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Could not start server: %v", err)
			os.Exit(1)
		}
	}()

	<-done
	slog.Info("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown Failed: %v", err)
		os.Exit(1)
	}
	slog.Info("Server Exited Properly")
}
