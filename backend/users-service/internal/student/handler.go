package student

import (
	"encoding/json"
	"errors"
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
		if identity.Role == commonauth.RoleStudent {
			item, err := h.service.Get(r.Context(), identity.StudentID)
			if err != nil {
				writeError(w, err)
				return
			}
			response.JSON(w, http.StatusOK, []Student{item})
			return
		}
		items, err := h.service.List(r.Context(), identity.TutorID)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.JSON(w, http.StatusOK, items)
	case http.MethodPost:
		if identity.Role != commonauth.RoleTutor {
			response.Error(w, http.StatusForbidden, commonauth.ErrForbidden.Error())
			return
		}
		var req CreateStudentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		item, err := h.service.Create(r.Context(), identity.TutorID, req)
		writeStudentResult(w, item, err, http.StatusCreated)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) ByID(w http.ResponseWriter, r *http.Request) {
	id, action, ok := parseStudentPath(r.URL.Path)
	if !ok {
		response.Error(w, http.StatusNotFound, "not found")
		return
	}
	identity, err := commonauth.Current(r.Context())
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	if identity.Role == commonauth.RoleStudent && id != identity.StudentID {
		response.Error(w, http.StatusForbidden, commonauth.ErrForbidden.Error())
		return
	}

	switch {
	case r.Method == http.MethodGet && action == "":
		item, err := h.service.Get(r.Context(), id)
		writeStudentResult(w, item, err, http.StatusOK)
	case r.Method == http.MethodPut && action == "":
		if identity.Role != commonauth.RoleTutor {
			response.Error(w, http.StatusForbidden, commonauth.ErrForbidden.Error())
			return
		}
		var req UpdateStudentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		item, err := h.service.Update(r.Context(), id, req)
		writeStudentResult(w, item, err, http.StatusOK)
	case r.Method == http.MethodDelete && action == "":
		if identity.Role != commonauth.RoleTutor {
			response.Error(w, http.StatusForbidden, commonauth.ErrForbidden.Error())
			return
		}
		err := h.service.Delete(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case r.Method == http.MethodPost && action == "accept":
		if identity.Role != commonauth.RoleTutor {
			response.Error(w, http.StatusForbidden, commonauth.ErrForbidden.Error())
			return
		}
		item, err := h.service.Accept(r.Context(), id)
		writeStudentResult(w, item, err, http.StatusOK)
	case r.Method == http.MethodPost && action == "archive":
		if identity.Role != commonauth.RoleTutor {
			response.Error(w, http.StatusForbidden, commonauth.ErrForbidden.Error())
			return
		}
		item, err := h.service.Archive(r.Context(), id)
		writeStudentResult(w, item, err, http.StatusOK)
	case r.Method == http.MethodPost && action == "notes":
		if identity.Role != commonauth.RoleTutor {
			response.Error(w, http.StatusForbidden, commonauth.ErrForbidden.Error())
			return
		}
		var req NotesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		item, err := h.service.UpdateNotes(r.Context(), id, req.Notes)
		writeStudentResult(w, item, err, http.StatusOK)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func parseStudentPath(path string) (int64, string, bool) {
	raw := strings.TrimPrefix(path, "/api/students/")
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

func writeStudentResult(w http.ResponseWriter, item Student, err error, status int) {
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
