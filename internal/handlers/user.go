package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LeonovDS/review-manager/internal/model"
)

// User contains dependencies for /user handlers.
type User struct {
	RG reviewGetter
}

type reviewGetter interface {
	Get(ctx context.Context, uID string) (model.ReviewReport, error)
}

// GetReview - GET /users/getReview - gets list of pull requests reviewed by one user.
func (u *User) GetReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uID := r.URL.Query().Get("user_id")

	report, err := u.RG.Get(ctx, uID)
	if err != nil {
		handleError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(report)
	if err != nil {
		slog.Error("Failed to write response", "err", err)
		return
	}
}
