package user

import (
	"context"

	"github.com/LeonovDS/review-manager/internal/model"
)

// StatusUpdater provides use case for updating user status.
type StatusUpdater struct {
	User userRepo
}

type userRepo interface {
	Get(ctx context.Context, id string) (model.TeamMember, error)
	SetIsActive(ctx context.Context, uID string, isActive bool) error
}

// SetIsActive updates isActive field of user.
func (r *StatusUpdater) SetIsActive(
	ctx context.Context, uID string, isActive bool,
) (model.TeamMember, error) {
	if len(uID) == 0 {
		return model.TeamMember{}, model.ErrBadRequest
	}

	err := r.User.SetIsActive(ctx, uID, isActive)
	if err != nil {
		return model.TeamMember{}, err
	}

	return r.User.Get(ctx, uID)
}
