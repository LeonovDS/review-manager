package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LeonovDS/review-manager/internal/model"
	"github.com/LeonovDS/review-manager/internal/usecase/user"
)

// UserHandler contains dependencies for /user handlers.
type UserHandler struct {
	reviews  *user.ReviewGetter
	isActive *user.StatusUpdater
}

// NewUserHandler creates new UserHandler.
func NewUserHandler(
	reviews *user.ReviewGetter,
	isActive *user.StatusUpdater,
) UserHandler {
	return UserHandler{
		reviews:  reviews,
		isActive: isActive,
	}
}

// GetReview - GET /users/getReview - gets list of pull requests reviewed by one user.
func (u *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uID := r.URL.Query().Get("user_id")

	report, err := u.reviews.Get(ctx, uID)
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

// SetIsActive - Post /users/setIsActive - sets user's isActive status.
func (u *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req setIsActiveRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, model.ErrBadRequest)
		return
	}

	user, err := u.isActive.SetIsActive(ctx, req.UID, req.IsActive)
	if err != nil {
		handleError(w, err)
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		slog.Error("Failed to write response", "err", err)
		return
	}
}
