package usecase

import (
	"github.com/LeonovDS/review-manager/internal/model"
)

type pullRequestRepository interface {
	Create(prID, prName, author string) (model.PullRequest, error)
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
