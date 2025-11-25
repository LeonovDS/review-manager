// Package pullrequest provides use cases for working with pull requests.
package pullrequest

import (
	"context"
	"math/rand"

	"github.com/LeonovDS/review-manager/internal/database"
	"github.com/LeonovDS/review-manager/internal/model"
)

// Creator provides use case for creating pull request.
type Creator struct {
	TX   database.TransactionManager
	PR   prCreatorRepo
	User userRepo
}

type prCreatorRepo interface {
	Create(ctx context.Context, prID, prName, authorID string) (model.PullRequest, error)
	AssignReviewers(ctx context.Context, prID string, reviewers []string) error
}

type userRepo interface {
	Get(ctx context.Context, id string) (model.User, error)
	GetActiveTeamMembers(ctx context.Context, user model.User) ([]string, error)
}

const maxReviewers int = 2

// Create validates request and saves pull request into repository.
func (u *Creator) Create(ctx context.Context, id, name, author string) (model.PullRequest, error) {
	err := validatePR(id, name, author)
	if err != nil {
		return model.PullRequest{}, err
	}

	var pr model.PullRequest
	err = u.TX.WithTransaction(ctx, func(context.Context) error {
		authorUser, err := u.User.Get(ctx, author)
		if err != nil {
			return err
		}

		teamMembers, err := u.User.GetActiveTeamMembers(ctx, authorUser)
		if err != nil {
			return err
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

		pr, err = u.PR.Create(ctx, id, name, author)
		if err != nil {
			return err
		}

		err = u.PR.AssignReviewers(ctx, pr.ID, reviewers)
		if err != nil {
			return err
		}

		pr.Reviewers = reviewers
		return nil
	})
	if err != nil {
		return model.PullRequest{}, err
	}

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
