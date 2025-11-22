// Package main starts review manager server.
package main

import (
	"context"
	"log/slog"

	"github.com/LeonovDS/review-manager/internal/database"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := database.Connect(ctx)
	if err != nil {
		slog.Error("Unable to perform migrations", "err", err)
		return
	}
}
