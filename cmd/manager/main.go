// Package main starts review manager server.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/LeonovDS/review-manager/internal/database"
	"github.com/LeonovDS/review-manager/internal/handlers"
	"github.com/LeonovDS/review-manager/internal/repository"
	pullrequest "github.com/LeonovDS/review-manager/internal/usecase/pull_request"
	"github.com/LeonovDS/review-manager/internal/usecase/team"
	"github.com/LeonovDS/review-manager/internal/usecase/user"
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

	teamRepo := repository.Team{Pool: pool}
	userRepo := repository.User{Pool: pool}
	prRepo := repository.PullRequest{Pool: pool}
	teamHandler := handlers.NewTeamHandler(
		&team.Adder{Team: &teamRepo, User: &userRepo},
		&team.Getter{Team: &teamRepo, User: &userRepo},
	)
	prHandler := handlers.NewPullRequestHandler(
		&pullrequest.Creator{PR: &prRepo, User: &userRepo},
		&pullrequest.Merger{PR: &prRepo},
		&pullrequest.Reassigner{PR: &prRepo, User: &userRepo},
	)
	userHandler := handlers.NewUserHandler(
		&user.ReviewGetter{PR: &prRepo},
		&user.StatusUpdater{User: &userRepo},
	)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /team/add", teamHandler.Add)
	mux.HandleFunc("GET /team/get", teamHandler.Get)
	mux.HandleFunc("POST /pullRequest/create", prHandler.Create)
	mux.HandleFunc("POST /pullRequest/merge", prHandler.Merge)
	mux.HandleFunc("POST /pullRequest/reassign", prHandler.Reassign)
	mux.HandleFunc("GET /users/getReview", userHandler.GetReview)
	mux.HandleFunc("POST /users/setIsActive", userHandler.SetIsActive)

	var server http.Server
	server.Addr = ":8080"
	server.Handler = mux
	server.ReadHeaderTimeout = 1 * time.Second

	err = server.ListenAndServe()
	if err != nil {
		slog.Error("Failed to start server", slog.Any("err", err))
		return
	}
}
