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
func (r *Team) Add(ctx context.Context, team model.Team) (model.Team, error) {
	var name string
	err := r.Pool.QueryRow(ctx, `
		INSERT INTO Team (name) 
		VALUES ($1) 
		ON CONFLICT (name) DO NOTHING 
		RETURNING name;
	`, team.TeamName).Scan(&name)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.Team{}, fmt.Errorf("%s %w", team.TeamName, model.ErrTeamExists)
	} else if err != nil {
		return model.Team{}, err
	}
	return model.Team{TeamName: name, Members: []model.User{}}, nil
}

// Get searches database for team with given name.
func (r *Team) Get(ctx context.Context, name string) (model.Team, error) {
	var dbName string
	err := r.Pool.QueryRow(ctx, `
		SELECT (name) 
		FROM Team 
		WHERE name=$1;
	`, name).Scan(&dbName)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.Team{}, fmt.Errorf("%s %w", name, model.ErrNotFound)
	} else if err != nil {
		return model.Team{}, err
	}

	return model.Team{TeamName: name, Members: []model.User{}}, nil
}
