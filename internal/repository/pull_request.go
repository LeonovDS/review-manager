package repository

import (
	"context"
	"errors"
	"time"

	"github.com/LeonovDS/review-manager/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PullRequest is database repository for working with pull requests.
type PullRequest struct {
	Pool *pgxpool.Pool
}

// Create creates pull request and returns it with created_at set.
func (r *PullRequest) Create(prID, prName, author string) (model.PullRequest, error) {
	var pr model.PullRequest
	err := r.Pool.QueryRow(context.TODO(), `
		INSERT INTO PullRequest (pull_request_id, pull_request_name, author_id, status)	
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (pull_request_id) DO NOTHING
		RETURNING pull_request_id, pull_request_name, author_id, status, created_at;
	`, prID, prName, author, "OPEN").Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.PullRequest{}, model.ErrPRExists
	} else if err != nil {
		return model.PullRequest{}, err
	}
	return pr, nil
}

// Merge updates pull request status.
// Cannot separate cases when PR is not found or not updated, so needs additional checks on call side.
func (r *PullRequest) Merge(id string) error {
	tag, err := r.Pool.Exec(context.TODO(), `
		UPDATE PullRequest
		SET status = 'MERGED', merged_at = NOW()
		WHERE pull_request_id = $1 
		AND status = 'OPEN';
	`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

// Get acquires pull request from repository.
func (r *PullRequest) Get(id string) (model.PullRequest, error) {
	var pr model.PullRequest
	var mergedAt time.Time
	err := r.Pool.QueryRow(context.TODO(), `
		SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		FROM PullRequest
		WHERE pull_request_id = $1;
	`, id).Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &mergedAt)
	pr.MergedAt = &mergedAt
	if errors.Is(err, pgx.ErrNoRows) {
		return model.PullRequest{}, model.ErrNotFound
	} else if err != nil {
		return model.PullRequest{}, err
	}

	return pr, nil
}
