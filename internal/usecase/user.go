package usecase

import (
	"context"

	"github.com/LeonovDS/review-manager/internal/model"
)

// GetReviews provides use case for forming report about user's reviews.
type GetReviews struct {
	PR reviewGetter
}

type reviewGetter interface {
	GetByUserID(ctx context.Context, uID string) ([]model.PullRequestShort, error)
}

// Get forms a report about user's reviews.
func (r *GetReviews) Get(ctx context.Context, uID string) (model.ReviewReport, error) {
	if len(uID) == 0 {
		return model.ReviewReport{}, model.ErrBadRequest
	}

	prs, err := r.PR.GetByUserID(ctx, uID)
	if err != nil {
		return model.ReviewReport{}, err
	}

	return model.ReviewReport{
		UserID:       uID,
		PullRequests: prs,
	}, nil
}
