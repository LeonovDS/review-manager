package usecase

import (
	"errors"

	"github.com/LeonovDS/review-manager/internal/model"
)

type pullRequestRepository interface {
	Create(prID, prName, author string) (model.PullRequest, error)
	Merge(id string) error
	Get(id string) (model.PullRequest, error)
}

type userRepositoty interface {
	UserExist(id string) error
}

// CreatePR provides use case for creating pull request.
type CreatePR struct {
	PR pullRequestRepository
	U  userRepositoty
}

// Create validates request and saves pull request into repository.
func (u *CreatePR) Create(id, name, author string) (model.PullRequest, error) {
	err := validatePR(id, name, author)
	if err != nil {
		return model.PullRequest{}, err
	}

	err = u.U.UserExist(author)
	if err != nil {
		return model.PullRequest{}, err
	}

	pr, err := u.PR.Create(id, name, author)
	if err != nil {
		return model.PullRequest{}, err
	}
	return pr, nil
}

func validatePR(id, name, author string) error {
	if len(id) == 0 || len(name) == 0 || len(author) == 0 {
		return model.ErrBadRequest
	}
	return nil
}

// MergePR provides use case for merging pull request.
type MergePR struct {
	PR pullRequestRepository
}

// Merge marks pull request as merged and returns it.
func (u *MergePR) Merge(id string) (model.PullRequest, error) {
	if len(id) == 0 {
		return model.PullRequest{}, model.ErrBadRequest
	}

	err := u.PR.Merge(id)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return model.PullRequest{}, err
	}

	pr, err := u.PR.Get(id)
	if err != nil {
		return model.PullRequest{}, err
	}
	return pr, nil
}
