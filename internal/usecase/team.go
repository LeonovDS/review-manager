// Package usecase provides domain logic.
package usecase

import "github.com/LeonovDS/review-manager/internal/model"

type teamRepository interface {
	Add(team model.Team) (model.Team, error)
	Get(name string) (model.Team, error)
}

// AddTeam provides use case for creating a new team.
type AddTeam struct {
	Repository teamRepository
}

// Add validates team and stores it into repository.
func (u *AddTeam) Add(team model.Team) (model.Team, error) {
	if len(team.TeamName) == 0 || len(team.Members) == 0 {
		return model.Team{}, model.ErrBadRequest
	}
	return u.Repository.Add(team)
}

// GetTeam provides use case for getting team from repository.
type GetTeam struct {
	Repository teamRepository
}

// Get queries repository for a team with given name.
func (u *GetTeam) Get(name string) (model.Team, error) {
	return u.Repository.Get(name)
}
