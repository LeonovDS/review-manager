// Package repository provides interactions with database.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/LeonovDS/review-manager/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Team stores repository dependencies.
type Team struct {
	Pool *pgxpool.Pool
}

// Add saves team to database or returns error, if team exists.
func (r *Team) Add(team model.Team) (model.Team, error) {
	var name string
	err := r.Pool.QueryRow(context.TODO(), `
		INSERT INTO Team (name) 
		VALUES ($1) 
		ON CONFLICT (name) DO NOTHING 
		RETURNING name
	`, team.TeamName).Scan(&name)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.Team{}, fmt.Errorf("%s %w", team.TeamName, model.ErrTeamExists)
	} else if err != nil {
		return model.Team{}, err
	}
	return model.Team{TeamName: name, Members: []model.TeamMember{}}, nil
}
