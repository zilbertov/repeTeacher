package notification

import (
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

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	identity, err := commonauth.Current(r.Context())
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	var items []Notification
	if identity.Role == commonauth.RoleStudent {
		items, err = h.service.ListForStudent(r.Context(), identity.StudentID)
	} else {
		items, err = h.service.List(r.Context(), identity.TutorID)
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *Handler) ByID(w http.ResponseWriter, r *http.Request) {
	id, action, ok := parseNotificationPath(r.URL.Path)
	if !ok {
		response.Error(w, http.StatusNotFound, "not found")
		return
	}

	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if _, err := commonauth.Current(r.Context()); err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	var item Notification
	var err error
	switch action {
	case "read":
		item, err = h.service.MarkRead(r.Context(), id)
	case "approve":
		item, err = h.service.Approve(r.Context(), id)
	case "reject":
		item, err = h.service.Reject(r.Context(), id)
	default:
		response.Error(w, http.StatusNotFound, "not found")
		return
	}

	if err != nil {
		writeError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func parseNotificationPath(path string) (int64, string, bool) {
	raw := strings.TrimPrefix(path, "/api/notifications/")
	raw = strings.Trim(raw, "/")
	parts := strings.Split(raw, "/")
	if len(parts) < 2 {
		return 0, "", false
	}
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, "", false
	}
	return id, parts[1], true
}

func parseIDQuery(r *http.Request, key string) (int64, bool) {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	return id, err == nil
}

func writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrNotFound) {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	response.Error(w, http.StatusInternalServerError, err.Error())
}
