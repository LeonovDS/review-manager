package pullrequest

import (
	"context"
	"errors"

	"github.com/LeonovDS/review-manager/internal/model"
)

// Merger provides use case for merging pull request.
type Merger struct {
	PR prMergerRepo
}

type prMergerRepo interface {
	Get(ctx context.Context, id string) (model.PullRequest, error)
	Merge(ctx context.Context, id string) error
}

// Merge marks pull request as merged and returns it.
func (u *Merger) Merge(ctx context.Context, id string) (model.PullRequest, error) {
	if len(id) == 0 {
		return model.PullRequest{}, model.ErrBadRequest
	}

	err := u.PR.Merge(ctx, id)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return model.PullRequest{}, err
	}

	pr, err := u.PR.Get(ctx, id)
	if err != nil {
		return model.PullRequest{}, err
	}
	return pr, nil
}
