// Package usecase provides domain logic.
package usecase

import "github.com/LeonovDS/review-manager/internal/model"

type teamRepository interface {
	Add(team model.Team) (model.Team, error)
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
