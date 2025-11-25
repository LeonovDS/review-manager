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

// SetIsActive provides use case for updating user status.
type SetIsActive struct {
	IAS isActiveSetter
	U   userRepositotyPR
}

type isActiveSetter interface {
	SetIsActive(ctx context.Context, uID string, isActive bool) error
}

// SetIsActive sets isActive field of user.
func (r *SetIsActive) SetIsActive(
	ctx context.Context, uID string, isActive bool,
) (model.TeamMember, error) {
	if len(uID) == 0 {
		return model.TeamMember{}, model.ErrBadRequest
	}

	err := r.IAS.SetIsActive(ctx, uID, isActive)
	if err != nil {
		return model.TeamMember{}, err
	}

	return r.U.GetUser(ctx, uID)
}
