package chat

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
		var items []Chat
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
		var req CreateChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		if identity.Role == commonauth.RoleStudent {
			req.StudentID = identity.StudentID
			if req.TutorID == 0 {
				req.TutorID = identity.TutorID
			}
		} else {
			req.TutorID = 0
		}
		item, err := h.service.Create(r.Context(), identity.TutorID, req)
		if err != nil {
			writeError(w, err)
			return
		}
		response.JSON(w, http.StatusCreated, item)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) Messages(w http.ResponseWriter, r *http.Request) {
	chatID, ok := parseChatID(r.URL.Path)
	if !ok {
		response.Error(w, http.StatusNotFound, "not found")
		return
	}
	identity, err := commonauth.Current(r.Context())
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		items, err := h.service.ListMessages(r.Context(), chatID)
		if err != nil {
			writeError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, items)
	case http.MethodPost:
		var req SendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		req.SenderType = identity.Role
		item, err := h.service.SendMessage(r.Context(), chatID, req)
		if err != nil {
			writeError(w, err)
			return
		}
		response.JSON(w, http.StatusCreated, item)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func parseChatID(path string) (int64, bool) {
	raw := strings.TrimPrefix(path, "/api/chats/")
	raw = strings.TrimSuffix(raw, "/messages")
	raw = strings.Trim(raw, "/")
	id, err := strconv.ParseInt(raw, 10, 64)
	return id, err == nil
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
	switch {
	case errors.Is(err, ErrBadRequest):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, err.Error())
	}
}
