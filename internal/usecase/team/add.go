// Package team provides use cases for working with teams.
package team

import (
	"context"

	"github.com/LeonovDS/review-manager/internal/database"
	"github.com/LeonovDS/review-manager/internal/model"
)

// Adder provides use case for creating a new team.
type Adder struct {
	TX   database.TransactionManager
	Team teamAdderRepo
	User userAdderRepo
}

type teamAdderRepo interface {
	Add(ctx context.Context, team model.Team) (model.Team, error)
}

type userAdderRepo interface {
	Add(ctx context.Context, team model.Team) error
}

// Add validates team and stores it into repository.
func (u *Adder) Add(ctx context.Context, team model.Team) (model.Team, error) {
	err := validate(team)
	if err != nil {
		return model.Team{}, err
	}

	err = u.TX.WithTransaction(ctx, func(context.Context) error {
		_, err = u.Team.Add(ctx, team)
		if err != nil {
			return err
		}

		err = u.User.Add(ctx, team)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.Team{}, err
	}

	return team, nil
}

func validate(team model.Team) error {
	if len(team.TeamName) == 0 {
		return model.ErrBadRequest
	}
	if len(team.Members) == 0 {
		return model.ErrBadRequest
	}

	for _, m := range team.Members {
		if len(m.UserID) == 0 {
			return model.ErrBadRequest
		}
		if len(m.Username) == 0 {
			return model.ErrBadRequest
		}
	}
	return nil
}
