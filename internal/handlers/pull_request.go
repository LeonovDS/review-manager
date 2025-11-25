package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LeonovDS/review-manager/internal/model"
)

// PullRequest provides handlers for pull request operations.
type PullRequest struct {
	CreateUC   createPR
	MergeUC    mergePR
	ReassignUC reassignPR
}

type createPR interface {
	Create(ctx context.Context, id, name, author string) (model.PullRequest, error)
}

type mergePR interface {
	Merge(ctx context.Context, id string) (model.PullRequest, error)
}

type createPRRequest struct {
	ID     string `json:"pull_request_id"`
	Name   string `json:"pull_request_name"`
	Author string `json:"author_id"`
}

// CreatePR - POST /pullRequest/create - creates a new pull request or returns error, if it exists.
func (h *PullRequest) CreatePR(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createPRRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, model.ErrBadRequest)
		return
	}

	pr, err := h.CreateUC.Create(ctx, req.ID, req.Name, req.Author)
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

// MergePR - POST /pullRequest/merge - merges a pull request (idempotent).
func (h *PullRequest) MergePR(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req mergePRRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, model.ErrBadRequest)
		return
	}

	pr, err := h.MergeUC.Merge(ctx, req.ID)
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

type reassignPR interface {
	Reassign(ctx context.Context, r model.Reviewer) (model.Reviewer, error)
}

// ReassignPR - POST /pullRequest/reassign - reassigns a pull request to other team member if it is possible.
func (h *PullRequest) ReassignPR(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var rev model.Reviewer
	err := json.NewDecoder(r.Body).Decode(&rev)
	if err != nil {
		handleError(w, model.ErrBadRequest)
		return
	}

	rev, err = h.ReassignUC.Reassign(ctx, rev)
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
