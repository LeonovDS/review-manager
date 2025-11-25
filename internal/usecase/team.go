// Package usecase provides domain logic.
package usecase

import (
	"context"

	"github.com/LeonovDS/review-manager/internal/model"
)

type teamRepository interface {
	Add(ctx context.Context, team model.Team) (model.Team, error)
	Get(ctx context.Context, name string) (model.Team, error)
}

type userRepository interface {
	Add(ctx context.Context, team model.Team) error
	GetByTeam(ctx context.Context, name string) ([]model.TeamMember, error)
}

// AddTeam provides use case for creating a new team.
type AddTeam struct {
	Team teamRepository
	User userRepository
}

// Add validates team and stores it into repository.
func (u *AddTeam) Add(ctx context.Context, team model.Team) (model.Team, error) {
	err := validate(team)
	if err != nil {
		return model.Team{}, err
	}

	_, err = u.Team.Add(ctx, team)
	if err != nil {
		return model.Team{}, err
	}

	err = u.User.Add(ctx, team)
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

// GetTeam provides use case for getting team from repository.
type GetTeam struct {
	Team teamRepository
	User userRepository
}

// Get queries repository for a team with given name.
func (u *GetTeam) Get(ctx context.Context, name string) (model.Team, error) {
	if len(name) == 0 {
		return model.Team{}, model.ErrBadRequest
	}

	team, err := u.Team.Get(ctx, name)
	if err != nil {
		return model.Team{}, err
	}

	team.Members, err = u.User.GetByTeam(ctx, name)
	if err != nil {
		return model.Team{}, err
	}
	return team, nil
}
