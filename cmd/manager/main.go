// Package main starts review manager server.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/LeonovDS/review-manager/internal/database"
	"github.com/LeonovDS/review-manager/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Unable to load .env file, using system environment", slog.Any("err", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := database.Connect(ctx)
	if err != nil {
		slog.Error("Unable to perform migrations", slog.Any("err", err))
		return
	}
	defer pool.Close()

	slog.Info("Database connection created")

	var server http.Server
	server.Addr = ":8080"
	server.Handler = handlers.NewRouter(pool)
	server.ReadHeaderTimeout = 1 * time.Second

	err = server.ListenAndServe()
	if err != nil {
		slog.Error("Failed to start server", slog.Any("err", err))
		return
	}
}
