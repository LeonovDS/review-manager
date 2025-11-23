package model

import (
	"errors"
)

// ErrBadRequest is used when input is malformed or incorrect.
var ErrBadRequest = errors.New("bad request")

// ErrTeamExists is used when team already exists in database.
var ErrTeamExists = errors.New("already exists")

// ErrNotFound is used when entity is not found in database.
var ErrNotFound = errors.New("not found")
