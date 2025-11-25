package model

// ReviewReport summarizes pull requests reviewed by a specific user.
type ReviewReport struct {
	UserID       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}

// PullRequestShort - short form of PullRequest for some responses.
type PullRequestShort struct {
	ID     string `json:"pull_request_id"`
	Name   string `json:"pull_request_name"`
	Author string `json:"author_id"`
	Status string `json:"status"`
}
