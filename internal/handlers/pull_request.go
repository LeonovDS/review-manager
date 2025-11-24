package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LeonovDS/review-manager/internal/model"
)

// PullRequest provides handlers for pull request operations.
type PullRequest struct {
	CreatePRUC createPR
}

type createPR interface {
	Create(id, name, author string) (model.PullRequest, error)
}

type createPRRequest struct {
	ID     string `json:"pull_request_id"`
	Name   string `json:"pull_request_name"`
	Author string `json:"author_id"`
}

// CreatePR - POST /pullRequest/create - creates a new pull request or returns error, if it exists.
func (h *PullRequest) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req createPRRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, model.ErrBadRequest)
		return
	}

	pr, err := h.CreatePRUC.Create(req.ID, req.Name, req.Author)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(pr)
	if err != nil {
		slog.Error("Failed to write response", "err", err)
		return
	}
}
