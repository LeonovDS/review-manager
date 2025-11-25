package team_test

import (
	"testing"

	"github.com/LeonovDS/review-manager/internal/model"
	"github.com/LeonovDS/review-manager/internal/usecase/team"
	"github.com/stretchr/testify/assert"
)

func TestTeamGet_Validate(t *testing.T) {
	var u team.Getter

	team, err := u.Get(t.Context(), "")

	assert.Equal(t, noTeam, team)
	assert.ErrorIs(t, err, model.ErrBadRequest)
}

type teamGetTestCase struct {
	testName     string
	prepareMocks func(tR *teamMockRepo, uR *userMockRepo)
	teamName     string
	expected     model.Team
	expectedErr  error
}

func TestTeamGet(t *testing.T) {
	tests := []teamGetTestCase{
		{
			testName: "Happy path",
			prepareMocks: func(tR *teamMockRepo, uR *userMockRepo) {
				_ = tR.On("Get", "team1").Return(
					model.Team{TeamName: "team1", Members: []model.User{}}, nil)
				_ = uR.On("GetByTeam", "team1").Return(sampleTeam.Members, nil)
			},
			teamName:    "team1",
			expected:    sampleTeam,
			expectedErr: nil,
		},
		{
			testName: "Not Found",
			prepareMocks: func(tR *teamMockRepo, uR *userMockRepo) {
				_ = tR.On("Get", "team1").Return(noTeam, model.ErrNotFound)
				_ = uR.On("GetByTeam", "team1").Return([]model.User{}, nil)
			},
			teamName:    "team1",
			expected:    noTeam,
			expectedErr: model.ErrNotFound,
		},
		{
			testName: "Empty members",
			prepareMocks: func(tR *teamMockRepo, uR *userMockRepo) {
				_ = tR.On("Get", "team1").Return(
					model.Team{TeamName: "team1", Members: []model.User{}}, nil)
				_ = uR.On("GetByTeam", "team1").Return([]model.User{}, nil)
			},
			teamName:    "team1",
			expected:    model.Team{TeamName: "team1", Members: []model.User{}},
			expectedErr: nil,
		},
		{
			testName: "Internal error",
			prepareMocks: func(tR *teamMockRepo, uR *userMockRepo) {
				_ = tR.On("Get", "team1").Return(sampleTeam, nil)
				_ = uR.On("GetByTeam", "team1").Return([]model.User{}, errInternal)
			},
			teamName:    "team1",
			expected:    noTeam,
			expectedErr: errInternal,
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			teamRepo := new(teamMockRepo)
			userRepo := new(userMockRepo)
			test.prepareMocks(teamRepo, userRepo)
			u := team.Getter{Team: teamRepo, User: userRepo}
			team, err := u.Get(t.Context(), test.teamName)
			assert.Equal(t, test.expected, team)
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}
