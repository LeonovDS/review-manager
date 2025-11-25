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
	RG   reviewGetter
	SIAS setIsActiveSetter
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

type setIsActiveRequest struct {
	UID      string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type setIsActiveSetter interface {
	SetIsActive(ctx context.Context, uID string, isActive bool) (model.TeamMember, error)
}

// SetIsActive - Post /users/setIsActive - sets user's isActive status.
func (u *User) SetIsActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req setIsActiveRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, model.ErrBadRequest)
		return
	}

	user, err := u.SIAS.SetIsActive(ctx, req.UID, req.IsActive)
	if err != nil {
		handleError(w, err)
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		slog.Error("Failed to write response", "err", err)
		return
	}
}
