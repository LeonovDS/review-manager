package usecase_test

import (
	"errors"
	"testing"

	"github.com/LeonovDS/review-manager/internal/model"
	"github.com/LeonovDS/review-manager/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Add(team model.Team) (model.Team, error) {
	args := m.Called(team)
	return args.Get(0).(model.Team), args.Error(1)
}

func TestTeamAdd(t *testing.T) {
	assert := assert.New(t)
	repo := new(mockRepo)
	u := usecase.AddTeam{repo}

	var noTeam model.Team

	team1 := model.Team{
		TeamName: "team1",
		Members: []model.TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: true},
		},
	}

	call := repo.On("Add", team1).Return(team1, nil)
	team, err := u.Add(team1)
	assert.Equal(team1, team)
	assert.Nil(err, "Given correct team, when team not in database, then err should be nil")

	team, err = u.Add(model.Team{
		TeamName: "",
		Members:  []model.TeamMember{{UserID: "u1", Username: "Alice", IsActive: true}},
	})
	assert.Equal(noTeam, team, "Given ")
	assert.ErrorIs(err, model.ErrBadRequest)

	team, err = u.Add(model.Team{
		TeamName: "invalidTeam2",
		Members:  []model.TeamMember{},
	})
	assert.Equal(noTeam, team)
	assert.ErrorIs(err, model.ErrBadRequest)

	call.Unset()
	repo.On("Add", team1).Return(noTeam, model.ErrTeamExists)
	team, err = u.Add(team1)
	assert.Equal(noTeam, team)
	assert.ErrorIs(err, model.ErrTeamExists)

	call.Unset()
	errInternal := errors.New("internal error")
	repo.On("Add", team1).Return(noTeam, errInternal)
	team, err = u.Add(team1)
	assert.Equal(noTeam, team)
	assert.ErrorIs(err, errInternal)
}
