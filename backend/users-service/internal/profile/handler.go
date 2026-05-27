package profile

import (
	"encoding/json"
	"errors"
	"net/http"

	commonauth "github.com/zilbertov/repe-teacher-common/auth"
	"github.com/zilbertov/repe-teacher-common/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	identity, err := commonauth.RequireTutor(r.Context())
	if err != nil {
		writeAuthError(w, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		item, err := h.service.Get(r.Context(), identity.TutorID)
		writeResult(w, item, err, http.StatusOK)
	case http.MethodPut:
		var req UpdateProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		item, err := h.service.Update(r.Context(), identity.TutorID, req)
		writeResult(w, item, err, http.StatusOK)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) Password(w http.ResponseWriter, r *http.Request) {
	identity, err := commonauth.RequireTutor(r.Context())
	if err != nil {
		writeAuthError(w, err)
		return
	}

	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := h.service.ChangePassword(r.Context(), identity.TutorID, req); err != nil {
		writeError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": "password changed"})
}

func (h *Handler) Settings(w http.ResponseWriter, r *http.Request) {
	identity, err := commonauth.RequireTutor(r.Context())
	if err != nil {
		writeAuthError(w, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		settings, err := h.service.GetSettings(r.Context(), identity.TutorID)
		writeResult(w, settings, err, http.StatusOK)
	case http.MethodPut:
		var req NotificationSettings
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		settings, err := h.service.UpdateSettings(r.Context(), identity.TutorID, req)
		writeResult(w, settings, err, http.StatusOK)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func writeResult(w http.ResponseWriter, data any, err error, status int) {
	if err != nil {
		writeError(w, err)
		return
	}
	response.JSON(w, status, data)
}

func writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrBadRequest) {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Error(w, http.StatusInternalServerError, err.Error())
}

func writeAuthError(w http.ResponseWriter, err error) {
	if errors.Is(err, commonauth.ErrForbidden) {
		response.Error(w, http.StatusForbidden, err.Error())
		return
	}
	response.Error(w, http.StatusUnauthorized, err.Error())
}
