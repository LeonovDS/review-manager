// Package handlers contains http handlers.
package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LeonovDS/review-manager/internal/model"
)

// Team contains dependencies for /team handlers.
type Team struct {
	UC addTeamUsecase
}

type addTeamUsecase interface {
	Add(t model.Team) (model.Team, error)
}

// AddTeam - POST /team/add - adds team and users.
func (h *Team) AddTeam(w http.ResponseWriter, r *http.Request) {
	var team model.Team
	err := json.NewDecoder(r.Body).Decode(&team)
	slog.Info("/team/add", "team", team)
	if err != nil {
		handleError(w, model.ErrBadRequest)
		return
	}

	team, err = h.UC.Add(team)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(team)
	if err != nil {
		slog.Error("Failed to write response", "err", err)
		return
	}
}
