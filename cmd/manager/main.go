// Package main starts review manager server.
package main

import (
	"context"
	"log/slog"

	"github.com/LeonovDS/review-manager/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Unable to load .env file, using system environment", slog.Any("err", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err = database.Connect(ctx)
	if err != nil {
		slog.Error("Unable to perform migrations", slog.Any("err", err))
		return
	}

	slog.Info("DB connection created")
}
