package usecase

import (
	"context"
	"errors"
	"math/rand"

	"github.com/LeonovDS/review-manager/internal/model"
)

type pullRequestRepository interface {
	Create(ctx context.Context,
		prID, prName, author string) (model.PullRequest, error)
	AssignReviewers(ctx context.Context, prID string, reviewers []string) error
	Merge(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (model.PullRequest, error)
}

type userRepositoty interface {
	GetTeam(ctx context.Context, id string) (string, error)
	GetActiveTeamMembers(ctx context.Context, userID, teamID string) ([]string, error)
}

// CreatePR provides use case for creating pull request.
type CreatePR struct {
	PR pullRequestRepository
	U  userRepositoty
}

const maxReviewers int = 2

// Create validates request and saves pull request into repository.
func (u *CreatePR) Create(ctx context.Context, id, name, author string) (model.PullRequest, error) {
	err := validatePR(id, name, author)
	if err != nil {
		return model.PullRequest{}, err
	}

	teamID, err := u.U.GetTeam(ctx, author)
	if err != nil {
		return model.PullRequest{}, err
	}

	teamMembers, err := u.U.GetActiveTeamMembers(ctx, author, teamID)
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
