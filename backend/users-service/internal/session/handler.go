package session

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zilbertov/repe-teacher-common/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}

	item, err := h.service.Login(r.Context(), req)
	if err != nil {
		writeError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrBadRequest):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrInvalidCredentials), errors.Is(err, ErrNotFound):
		response.Error(w, http.StatusUnauthorized, ErrInvalidCredentials.Error())
	default:
		response.Error(w, http.StatusInternalServerError, err.Error())
	}
}
