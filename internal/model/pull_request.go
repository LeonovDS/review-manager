package model

import "time"

// PullRequest represents a pull request stored in the system.
type PullRequest struct {
	ID        string     `json:"pull_request_id"`
	Name      string     `json:"pull_request_name"`
	AuthorID  string     `json:"author_id"`
	Status    string     `json:"status"`
	Reviewers []string   `json:"assigned_reviewers"`
	CreatedAt *time.Time `json:"createdAt"`
	MergedAt  *time.Time `json:"mergedAt,omitempty"`
}

// Reviewer identifies a reviewer assigned to a pull request.
type Reviewer struct {
	PRID string `json:"pull_request_id"`
	UID  string `json:"old_user_id"`
}
