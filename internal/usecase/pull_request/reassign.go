package pullrequest

import (
	"context"
	"math/rand"
	"slices"

	"github.com/LeonovDS/review-manager/internal/model"
)

// Reassigner provides use case for reassigning pull requests.
type Reassigner struct {
	PR   prReassignerRepo
	User userRepo
}

type prReassignerRepo interface {
	Get(ctx context.Context, id string) (model.PullRequest, error)
	UpdateReviewer(ctx context.Context, prID, oUID, nUID string) error
}

// Reassign checks if user is actual reviewer of pull request and finds active team member who can review PR instead.
func (u *Reassigner) Reassign(ctx context.Context, r model.Reviewer) (model.Reviewer, error) {
	if len(r.PRID) == 0 || len(r.UID) == 0 {
		return model.Reviewer{}, model.ErrBadRequest
	}

	user, err := u.User.Get(ctx, r.UID)
	if err != nil {
		return model.Reviewer{}, err
	}

	pr, err := u.PR.Get(ctx, r.PRID)
	if err != nil {
		return model.Reviewer{}, err
	}
	if pr.Status == "MERGED" {
		return model.Reviewer{}, model.ErrPRMerged
	}
	if !slices.Contains(pr.Reviewers, r.UID) {
		return model.Reviewer{}, model.ErrNotAssigned
	}

	team, err := u.User.GetActiveTeamMembers(ctx, user)
	if err != nil {
		return model.Reviewer{}, err
	}

	team = slices.DeleteFunc(team, func(it string) bool {
		return slices.Contains(pr.Reviewers, it) || it == pr.AuthorID
	})

	if len(team) == 0 {
		return model.Reviewer{}, model.ErrNoCandidate
	}

	// #nosec G404
	newID := team[rand.Intn(len(team))]
	err = u.PR.UpdateReviewer(ctx, pr.ID, r.UID, newID)
	if err != nil {
		return model.Reviewer{}, err
	}

	return model.Reviewer{
		PRID: r.PRID,
		UID:  newID,
	}, nil
}
