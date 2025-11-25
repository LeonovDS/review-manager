// Package model provides types domain and API types.
package model

// Team represents a group of users participating in reviews.
type Team struct {
	TeamName string `json:"team_name"`
	Members  []User `json:"members"`
}

// User represents an application user and their team membership.
type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name,omitempty"`
}
