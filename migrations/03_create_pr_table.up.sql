CREATE TYPE Status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE IF NOT EXISTS PullRequest (
    pull_request_id TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES Users(user_id),
    status Status NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    merged_at TIMESTAMP
);
