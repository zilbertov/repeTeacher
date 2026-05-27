package lesson

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	commonauth "github.com/zilbertov/repe-teacher-common/auth"
	"github.com/zilbertov/repe-teacher-common/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	identity, err := commonauth.Current(r.Context())
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		var items []Lesson
		if identity.Role == commonauth.RoleStudent {
			items, err = h.service.ListByStudent(r.Context(), identity.StudentID)
		} else {
			items, err = h.service.List(r.Context(), identity.TutorID)
		}
		if err != nil {
			writeError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, items)
	case http.MethodPost:
		if identity.Role != commonauth.RoleTutor {
			response.Error(w, http.StatusForbidden, commonauth.ErrForbidden.Error())
			return
		}
		var req CreateLessonRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		req.TutorID = 0
		item, err := h.service.Create(r.Context(), identity.TutorID, req)
		writeLessonResult(w, item, err, http.StatusCreated)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) ByID(w http.ResponseWriter, r *http.Request) {
	id, action, ok := parseLessonPath(r.URL.Path)
	if !ok {
		response.Error(w, http.StatusNotFound, "not found")
		return
	}
	identity, err := commonauth.Current(r.Context())
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	switch {
	case r.Method == http.MethodGet && action == "":
		item, err := h.service.Get(r.Context(), id)
		writeLessonResult(w, item, err, http.StatusOK)
	case r.Method == http.MethodPut && action == "":
		if identity.Role != commonauth.RoleTutor {
			response.Error(w, http.StatusForbidden, commonauth.ErrForbidden.Error())
			return
		}
		var req UpdateLessonRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		item, err := h.service.Update(r.Context(), id, req)
		writeLessonResult(w, item, err, http.StatusOK)
	case r.Method == http.MethodPost && action == "reschedule":
		var req RescheduleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		req.SenderType = identity.Role
		item, err := h.service.Reschedule(r.Context(), id, req)
		writeLessonResult(w, item, err, http.StatusOK)
	case r.Method == http.MethodPost && action == "cancel":
		if identity.Role != commonauth.RoleTutor {
			response.Error(w, http.StatusForbidden, commonauth.ErrForbidden.Error())
			return
		}
		var req CancelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
			response.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		req.SenderType = identity.Role
		item, err := h.service.Cancel(r.Context(), id, req)
		writeLessonResult(w, item, err, http.StatusOK)
	case r.Method == http.MethodPost && action == "files":
		if identity.Role != commonauth.RoleTutor {
			response.Error(w, http.StatusForbidden, commonauth.ErrForbidden.Error())
			return
		}
		var req AddFileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		file, err := h.service.AddFile(r.Context(), id, req)
		if err != nil {
			writeError(w, err)
			return
		}
		response.JSON(w, http.StatusCreated, file)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func parseLessonPath(path string) (int64, string, bool) {
	raw := strings.TrimPrefix(path, "/api/lessons/")
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

func parseIDQuery(r *http.Request, key string) (int64, bool) {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	return id, err == nil
}

func writeLessonResult(w http.ResponseWriter, item Lesson, err error, status int) {
	if err != nil {
		writeError(w, err)
		return
	}
	response.JSON(w, status, item)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrBadRequest):
		response.Error(w, http.StatusBadRequest, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, err.Error())
	}
}
