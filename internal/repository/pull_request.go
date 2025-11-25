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
func (r *PullRequest) Create(
	ctx context.Context,
	prID, prName, author string,
) (model.PullRequest, error) {
	var pr model.PullRequest
	err := r.Pool.QueryRow(ctx, `
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

// AssignReviewers adds reviewers to pull requests WITHOUT ADDITIONAL CHECKS.
func (r *PullRequest) AssignReviewers(ctx context.Context, prID string, reviewers []string) error {
	var batch pgx.Batch
	query := `
		INSERT INTO UsersToPullRequests (pull_request_id, reviewer_id)
		VALUES ($1, $2);
	`
	for _, rID := range reviewers {
		batch.Queue(query, prID, rID)
	}
	br := r.Pool.SendBatch(ctx, &batch)
	defer func() { _ = br.Close() }()
	for range reviewers {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

// Merge updates pull request status.
// Cannot separate cases when PR is not found or not updated, so needs additional checks on call side.
func (r *PullRequest) Merge(ctx context.Context, id string) error {
	tag, err := r.Pool.Exec(ctx, `
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

// Get acquires pull request with reviewers from repository.
func (r *PullRequest) Get(ctx context.Context, id string) (model.PullRequest, error) {
	query := `
		SELECT
			pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status, pr.created_at, pr.merged_at,
			COALESCE(array_agg(rev.reviewer_id) FILTER (WHERE rev.reviewer_id IS NOT NULL), '{}') AS reviewers
		FROM PullRequest pr
		LEFT JOIN UsersToPullRequests rev 
			ON rev.pull_request_id = pr.pull_request_id
		WHERE pr.pull_request_id = $1
		GROUP BY pr.pull_request_id;
	`

	var pr model.PullRequest
	var mergedAt time.Time
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &mergedAt, &pr.Reviewers)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.PullRequest{}, model.ErrNotFound
	} else if err != nil {
		return model.PullRequest{}, err
	}

	pr.MergedAt = &mergedAt
	return pr, nil
}
