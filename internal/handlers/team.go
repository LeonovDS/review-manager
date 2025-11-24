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
	AddUC addTeamUsecase
	GetUC getTeamUsecase
}

type addTeamUsecase interface {
	Add(t model.Team) (model.Team, error)
}

type getTeamUsecase interface {
	Get(name string) (model.Team, error)
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

	team, err = h.AddUC.Add(team)
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

// GetTeam - GET /team/get - get team info.
func (h *Team) GetTeam(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("team_name")

	team, err := h.GetUC.Get(name)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(team)
	if err != nil {
		slog.Error("Failed to write response", "err", err)
		return
	}
}
