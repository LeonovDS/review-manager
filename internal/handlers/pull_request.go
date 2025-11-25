package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LeonovDS/review-manager/internal/model"
	"github.com/LeonovDS/review-manager/internal/usecase"
)

// PullRequestHandler provides handlers for pull request operations.
type PullRequestHandler struct {
	create   *usecase.CreatePR
	merge    *usecase.MergePR
	reassign *usecase.ReassignPR
}

// NewPullRequestHandler creates new PullRequestHandler.
func NewPullRequestHandler(
	create *usecase.CreatePR,
	merge *usecase.MergePR,
	reassign *usecase.ReassignPR,
) PullRequestHandler {
	return PullRequestHandler{
		create:   create,
		merge:    merge,
		reassign: reassign,
	}
}

type createPRRequest struct {
	ID     string `json:"pull_request_id"`
	Name   string `json:"pull_request_name"`
	Author string `json:"author_id"`
}

// Create - POST /pullRequest/create - creates a new pull request or returns error, if it exists.
func (h *PullRequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createPRRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, model.ErrBadRequest)
		return
	}

	pr, err := h.create.Create(ctx, req.ID, req.Name, req.Author)
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

type mergePRRequest struct {
	ID string `json:"pull_request_id"`
}

// Merge - POST /pullRequest/merge - merges a pull request (idempotent).
func (h *PullRequestHandler) Merge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req mergePRRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, model.ErrBadRequest)
		return
	}

	pr, err := h.merge.Merge(ctx, req.ID)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(pr)
	if err != nil {
		slog.Error("Failed to write response", "err", err)
		return
	}
}

// Reassign - POST /pullRequest/reassign - reassigns a pull request to other team member if it is possible.
func (h *PullRequestHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var rev model.Reviewer
	err := json.NewDecoder(r.Body).Decode(&rev)
	if err != nil {
		handleError(w, model.ErrBadRequest)
		return
	}

	rev, err = h.reassign.Reassign(ctx, rev)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(rev)
	if err != nil {
		slog.Error("Failed to write response", "err", err)
		return
	}
}
