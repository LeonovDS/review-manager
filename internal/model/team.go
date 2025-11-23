// Package model provides types used across layers.
package model

// Team used in API, data and domain layers.
type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

// TeamMember represents User only as part of Team type.
type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
