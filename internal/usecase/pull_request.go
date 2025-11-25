package usecase

import (
	"context"
	"errors"
	"math/rand"
	"slices"

	"github.com/LeonovDS/review-manager/internal/model"
)

type pullRequestRepository interface {
	Create(ctx context.Context,
		prID, prName, author string) (model.PullRequest, error)
	AssignReviewers(ctx context.Context, prID string, reviewers []string) error
	Merge(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (model.PullRequest, error)
	UpdateReviewer(ctx context.Context, prID, oUID, nUID string) error
}

type userRepositotyPR interface {
	GetActiveTeamMembers(ctx context.Context, userID, teamID string) ([]string, error)
	GetUser(ctx context.Context, userID string) (model.TeamMember, error)
}

// CreatePR provides use case for creating pull request.
type CreatePR struct {
	PR pullRequestRepository
	U  userRepositotyPR
}

const maxReviewers int = 2

// Create validates request and saves pull request into repository.
func (u *CreatePR) Create(ctx context.Context, id, name, author string) (model.PullRequest, error) {
	err := validatePR(id, name, author)
	if err != nil {
		return model.PullRequest{}, err
	}

	user, err := u.U.GetUser(ctx, author)
	if err != nil {
		return model.PullRequest{}, err
	}

	teamMembers, err := u.U.GetActiveTeamMembers(ctx, author, user.TeamName)
	if err != nil {
		return model.PullRequest{}, err
	}

	reviewers := make([]string, 0, maxReviewers)
	switch len(teamMembers) {
	case 0:
		break
	case 1:
		reviewers = append(reviewers, teamMembers[0])
	default:
		i1, i2 := randomPair(len(teamMembers))
		reviewers = append(reviewers, teamMembers[i1])
		reviewers = append(reviewers, teamMembers[i2])
	}

	pr, err := u.PR.Create(ctx, id, name, author)
	if err != nil {
		return model.PullRequest{}, err
	}

	err = u.PR.AssignReviewers(ctx, pr.ID, reviewers)
	if err != nil {
		return model.PullRequest{}, err
	}

	pr.Reviewers = reviewers
	return pr, nil
}

// #nosec G404 - there is no need for secure random
func randomPair(n int) (int, int) {
	a := rand.Intn(n)
	b := a
	for b == a {
		b = rand.Intn(n)
	}
	return a, b
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
func (u *MergePR) Merge(ctx context.Context, id string) (model.PullRequest, error) {
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

// ReassignPR provides use case for reassigning pull requests.
type ReassignPR struct {
	PR pullRequestRepository
	U  userRepositotyPR
}

// Reassign checks if user is actual reviewer of pull request and finds active team member who can review PR instead.
func (u *ReassignPR) Reassign(ctx context.Context, r model.Reviewer) (model.Reviewer, error) {
	if len(r.PRID) == 0 || len(r.UID) == 0 {
		return model.Reviewer{}, model.ErrBadRequest
	}

	user, err := u.U.GetUser(ctx, r.UID)
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

	team, err := u.U.GetActiveTeamMembers(ctx, r.UID, user.TeamName)
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
