package team

import (
	"context"

	"github.com/LeonovDS/review-manager/internal/model"
)

// Getter provides use case for getting team from repository.
type Getter struct {
	Team teamGetterRepository
	User userGetterRepository
}

type teamGetterRepository interface {
	Get(ctx context.Context, name string) (model.Team, error)
}

type userGetterRepository interface {
	GetByTeam(ctx context.Context, name string) ([]model.TeamMember, error)
}

// Get queries repository for a team with given name.
func (u *Getter) Get(ctx context.Context, name string) (model.Team, error) {
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
