package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LeonovDS/review-manager/internal/model"
)

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func handleError(w http.ResponseWriter, err error) {
	var data errorResponse
	data.Message = err.Error()
	var code int
	switch {
	case errors.Is(err, model.ErrBadRequest):
		data.Code = "BAD_REQUEST"
		code = http.StatusBadRequest
	case errors.Is(err, model.ErrTeamExists):
		data.Code = "TEAM_EXISTS"
		code = http.StatusBadRequest
	case errors.Is(err, model.ErrPRExists):
		data.Code = "PR_EXISTS"
		code = http.StatusConflict
	case errors.Is(err, model.ErrNotFound):
		data.Code = "NOT_FOUND"
		code = http.StatusNotFound
	default:
		data.Code = "INTERNAL_ERROR"
		code = http.StatusInternalServerError
	}
	errString, _ := json.Marshal(data)
	http.Error(w, string(errString), code)
}
