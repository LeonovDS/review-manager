package team_test

import (
	"context"
	"errors"
	"testing"

	"github.com/LeonovDS/review-manager/internal/model"
	"github.com/LeonovDS/review-manager/internal/usecase/team"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//nolint:gochecknoglobals
var (
	errInternal = errors.New("internal error")
	sampleTeam  = model.Team{
		TeamName: "team1",
		Members: []model.User{
			{UserID: "u1", Username: "Alice", IsActive: true, TeamName: ""},
			{UserID: "u2", Username: "Bob", IsActive: true, TeamName: ""},
		},
	}
	noTeam model.Team
)

type teamMockRepo struct {
	mock.Mock
}

func (m *teamMockRepo) Add(_ context.Context, team model.Team) (model.Team, error) {
	args := m.Called(team)
	return args.Get(0).(model.Team), args.Error(1)
}

func (m *teamMockRepo) Get(_ context.Context, name string) (model.Team, error) {
	args := m.Called(name)
	return args.Get(0).(model.Team), args.Error(1)
}

type userMockRepo struct {
	mock.Mock
}

func (m *userMockRepo) Add(_ context.Context, team model.Team) error {
	args := m.Called(team)
	return args.Error(0)
}

func (m *userMockRepo) GetByTeam(_ context.Context, name string) ([]model.User, error) {
	args := m.Called(name)
	return args.Get(0).([]model.User), args.Error(1)
}

func TestTeamAdd_Validation(t *testing.T) {
	var u team.Adder

	type testCase struct {
		testName string
		team     model.Team
	}

	tests := []testCase{
		{
			testName: "Empty TeamName",
			team: model.Team{TeamName: "", Members: []model.User{
				{UserID: "u1", Username: "Alice", IsActive: true, TeamName: ""},
			}},
		},
		{
			testName: "Empty Members",
			team:     model.Team{TeamName: "team1", Members: []model.User{}},
		},
		{
			testName: "Empty UserID",
			team: model.Team{TeamName: "team1", Members: []model.User{
				{UserID: "u1", Username: "Alice", IsActive: true, TeamName: ""},
				{UserID: "", Username: "Bob", IsActive: true, TeamName: ""},
			}},
		},
		{
			testName: "Empty Username",
			team: model.Team{TeamName: "team1", Members: []model.User{
				{UserID: "u1", Username: "Alice", IsActive: true, TeamName: ""},
				{UserID: "u2", Username: "", IsActive: true, TeamName: ""},
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			team, err := u.Add(t.Context(), test.team)
			assert.ErrorIs(t, err, model.ErrBadRequest, "Expected ErrBadRequest")
			assert.Equal(t, noTeam, team, "Expected default Team")
		})
	}
}

func TestTeamAdd(t *testing.T) {
	teamRepo := new(teamMockRepo)
	userRepo := new(userMockRepo)
	teamRepo.On("Add", sampleTeam).Return(sampleTeam, nil)
	userRepo.On("Add", mock.Anything).Return(nil)

	u := team.Adder{Team: teamRepo, User: userRepo}
	team, err := u.Add(t.Context(), sampleTeam)
	assert.Equal(t, sampleTeam, team)
	assert.NoError(t, err)
}

func TestTeamAdd_Errors(t *testing.T) {
	type testCase struct {
		testName     string
		prepareMocks func(tR *teamMockRepo, uR *userMockRepo)
		input        model.Team
		expectedErr  error
	}

	tests := []testCase{
		{
			testName: "TeamRepo internal error",
			prepareMocks: func(tR *teamMockRepo, uR *userMockRepo) {
				_ = tR.On("Add", sampleTeam).Return(noTeam, errInternal)
				_ = uR.On("Add", mock.Anything).Return(nil)
			},
			input:       sampleTeam,
			expectedErr: errInternal,
		},
		{
			testName: "UserRepo internal error",
			prepareMocks: func(tR *teamMockRepo, uR *userMockRepo) {
				_ = tR.On("Add", sampleTeam).Return(sampleTeam, nil)
				_ = uR.On("Add", mock.Anything).Return(errInternal)
			},
			input:       sampleTeam,
			expectedErr: errInternal,
		},
		{
			testName: "TeamExists error",
			prepareMocks: func(tR *teamMockRepo, uR *userMockRepo) {
				_ = tR.On("Add", sampleTeam).Return(noTeam, model.ErrTeamExists)
				_ = uR.On("Add", mock.Anything).Return(nil)
			},
			input:       sampleTeam,
			expectedErr: model.ErrTeamExists,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			teamRepo := new(teamMockRepo)
			userRepo := new(userMockRepo)
			test.prepareMocks(teamRepo, userRepo)
			u := team.Adder{Team: teamRepo, User: userRepo}
			team, err := u.Add(t.Context(), test.input)
			assert.Equal(t, noTeam, team)
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}
