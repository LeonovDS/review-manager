// Package handlers contains router and http handlers.
package handlers

import (
	"net/http"

	"github.com/LeonovDS/review-manager/internal/database"
	"github.com/LeonovDS/review-manager/internal/repository"
	pullrequest "github.com/LeonovDS/review-manager/internal/usecase/pull_request"
	"github.com/LeonovDS/review-manager/internal/usecase/team"
	"github.com/LeonovDS/review-manager/internal/usecase/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewRouter builds handlers from dependencies and combine them into router.
func NewRouter(pool *pgxpool.Pool) *http.ServeMux {
	tm := database.DBTransactionManager{Pool: pool}
	teamRepo := repository.Team{Pool: pool}
	userRepo := repository.User{Pool: pool}
	prRepo := repository.PullRequest{Pool: pool}
	teamHandler := NewTeamHandler(
		&team.Adder{TX: &tm, Team: &teamRepo, User: &userRepo},
		&team.Getter{Team: &teamRepo, User: &userRepo},
	)
	prHandler := NewPullRequestHandler(
		&pullrequest.Creator{TX: &tm, PR: &prRepo, User: &userRepo},
		&pullrequest.Merger{PR: &prRepo},
		&pullrequest.Reassigner{TX: &tm, PR: &prRepo, User: &userRepo},
	)
	userHandler := NewUserHandler(
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

	return mux
}
