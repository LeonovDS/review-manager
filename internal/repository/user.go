package repository

import (
	"context"
	"errors"

	"github.com/LeonovDS/review-manager/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// User stores user repository dependencies.
type User struct {
	Pool *pgxpool.Pool
}

// Add saves users to database.
func (r *User) Add(ctx context.Context, t model.Team) error {
	query := `
		INSERT INTO Users (user_id, username, is_active, team)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE 
		SET username = EXCLUDED.username, is_active = EXCLUDED.is_active, team = EXCLUDED.team;
	`

	var batch pgx.Batch
	for _, u := range t.Members {
		batch.Queue(query, u.UserID, u.Username, u.IsActive, t.TeamName)
	}
	br := r.Pool.SendBatch(ctx, &batch)
	defer func() { _ = br.Close() }()

	for range t.Members {
		_, err := br.Exec()
		if err != nil {
			return errors.New("internal error")
		}
	}
	return nil
}

// GetByTeam acquires team members from one team.
func (r *User) GetByTeam(ctx context.Context, teamName string) ([]model.TeamMember, error) {
	query := `
		SELECT user_id, username, is_active
		FROM Users 
		WHERE team = $1;
	`
	rows, err := r.Pool.Query(ctx, query, teamName)
	if err != nil {
		return []model.TeamMember{}, err
	}
	defer rows.Close()

	results := []model.TeamMember{}
	for rows.Next() {
		var member model.TeamMember
		err := rows.Scan(&member.UserID, &member.Username, &member.IsActive)
		if err != nil {
			return nil, err
		}
		results = append(results, member)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetTeam find user's team by userID and returns error if user is missing.
func (r *User) GetTeam(ctx context.Context, id string) (string, error) {
	query := `
		SELECT team 
		FROM Users 
		WHERE user_id = $1;
	`
	var team string
	err := r.Pool.QueryRow(ctx, query, id).Scan(&team)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", model.ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return team, nil
}

// GetActiveTeamMembers finds other active users from the same team.
func (r *User) GetActiveTeamMembers(ctx context.Context, userID, teamID string) ([]string, error) {
	query := `
		SELECT user_id 
		FROM Users 
		WHERE team = $1 
			AND user_id <> $2
			AND is_active;
	`
	var teams []string
	rows, err := r.Pool.Query(ctx, query, teamID, userID)
	if err != nil {
		return teams, err
	}

	for rows.Next() {
		var teamID string
		err := rows.Scan(&teamID)
		if err != nil {
			return []string{}, err
		}
		teams = append(teams, teamID)
	}

	err = rows.Err()
	if err != nil {
		return []string{}, err
	}
	return teams, nil
}
