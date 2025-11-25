package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LeonovDS/review-manager/internal/model"
	"github.com/LeonovDS/review-manager/internal/usecase/team"
)

// TeamHandler contains dependencies for /team handlers.
type TeamHandler struct {
	add *team.Adder
	get *team.Getter
}

// NewTeamHandler creates new TeamHandler.
func NewTeamHandler(add *team.Adder, get *team.Getter) TeamHandler {
	return TeamHandler{
		add: add,
		get: get,
	}
}

// Add - POST /team/add - adds team and users.
func (h *TeamHandler) Add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var team model.Team
	err := json.NewDecoder(r.Body).Decode(&team)
	slog.Info("/team/add", "team", team)
	if err != nil {
		handleError(w, model.ErrBadRequest)
		return
	}

	team, err = h.add.Add(ctx, team)
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

// Get - GET /team/get - get team info.
func (h *TeamHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := r.URL.Query().Get("team_name")

	team, err := h.get.Get(ctx, name)
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
