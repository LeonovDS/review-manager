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
func (r *User) Add(t model.Team) error {
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
	br := r.Pool.SendBatch(context.TODO(), &batch)
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
func (r *User) GetByTeam(teamName string) ([]model.TeamMember, error) {
	query := `
		SELECT user_id, username, is_active
		FROM Users 
		WHERE team = $1;
	`
	rows, err := r.Pool.Query(context.TODO(), query, teamName)
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
