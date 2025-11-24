package repository

import (
	"context"
	"errors"

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
