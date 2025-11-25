package model

import (
	"errors"
)

// ErrBadRequest is used when input is malformed or incorrect.
var ErrBadRequest = errors.New("bad request")

// ErrTeamExists is used when team already exists in database.
var ErrTeamExists = errors.New("already exists")

// ErrPRExists is used when pull request already exists in database.
var ErrPRExists = errors.New("already exists")

// ErrNotFound is used when entity is not found in database.
var ErrNotFound = errors.New("not found")

// ErrNoCandidate is used when there are no candidate to reassign review.
var ErrNoCandidate = errors.New("no candidate")

// ErrNotAssigned is used when user is not reviewer, but tries to reassign review.
var ErrNotAssigned = errors.New("not assigned")

// ErrPRMerged is used when attempted to reassign review of already merged pull request.
var ErrPRMerged = errors.New("pr merged")
