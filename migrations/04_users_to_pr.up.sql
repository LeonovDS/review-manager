CREATE TABLE IF NOT EXISTS UsersToPullRequests (
    pull_request_id TEXT NOT NULL REFERENCES PullRequest(pull_request_id),
    reviewer_id TEXT NOT NULL REFERENCES Users(user_id),
    PRIMARY KEY (pull_request_id, reviewer_id)
);
