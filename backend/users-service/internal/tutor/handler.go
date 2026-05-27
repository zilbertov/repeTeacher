package tutor

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/zilbertov/repe-teacher-common/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := h.service.List(r.Context())
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.JSON(w, http.StatusOK, items)
	case http.MethodPost:
		var req CreateTutorRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		item, err := h.service.Create(r.Context(), req)
		writeResult(w, item, err, http.StatusCreated)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) ByID(w http.ResponseWriter, r *http.Request) {
	id, action, ok := parseTutorPath(r.URL.Path)
	if !ok {
		response.Error(w, http.StatusNotFound, "not found")
		return
	}

	if r.Method == http.MethodGet && action == "" {
		item, err := h.service.Get(r.Context(), id)
		writeResult(w, item, err, http.StatusOK)
		return
	}

	if r.Method == http.MethodPost && action == "notes" {
		var req UpdateTutorNotesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		item, err := h.service.UpdateNotes(r.Context(), id, req)
		writeResult(w, item, err, http.StatusOK)
		return
	}

	if action != "" {
		response.Error(w, http.StatusNotFound, "not found")
		return
	}

	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
}

func parseTutorPath(path string) (int64, string, bool) {
	raw := strings.TrimPrefix(path, "/api/tutors/")
	raw = strings.Trim(raw, "/")
	if raw == "" {
		return 0, "", false
	}
	parts := strings.Split(raw, "/")
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, "", false
	}
	action := ""
	if len(parts) > 1 {
		action = parts[1]
	}
	return id, action, true
}

func writeResult(w http.ResponseWriter, item Tutor, err error, status int) {
	if err != nil {
		writeError(w, err)
		return
	}
	response.JSON(w, status, item)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrBadRequest):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, err.Error())
	}
}
