package chat

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	commonauth "github.com/zilbertov/repe-teacher-common/auth"
)

func TestHandlerListChats(t *testing.T) {
	handler := NewHandler(NewService(&fakeChatRepo{}))
	req := withTutor(httptest.NewRequest(http.MethodGet, "/api/chats", nil))
	rec := httptest.NewRecorder()

	handler.ListOrCreate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandlerListChatsByStudent(t *testing.T) {
	handler := NewHandler(NewService(&fakeChatRepo{}))
	req := withStudent(httptest.NewRequest(http.MethodGet, "/api/chats?student_id=4", nil))
	rec := httptest.NewRecorder()

	handler.ListOrCreate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandlerCreateChat(t *testing.T) {
	handler := NewHandler(NewService(&fakeChatRepo{}))
	body := bytes.NewBufferString(`{"student_id":4}`)
	req := withTutor(httptest.NewRequest(http.MethodPost, "/api/chats", body))
	rec := httptest.NewRecorder()

	handler.ListOrCreate(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandlerMessages(t *testing.T) {
	handler := NewHandler(NewService(&fakeChatRepo{}))

	sendReq := withStudent(httptest.NewRequest(http.MethodPost, "/api/chats/1/messages", bytes.NewBufferString(`{"sender_type":"student","text":"Привет"}`)))
	sendRec := httptest.NewRecorder()
	handler.Messages(sendRec, sendReq)
	if sendRec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", sendRec.Code, sendRec.Body.String())
	}

	listReq := withStudent(httptest.NewRequest(http.MethodGet, "/api/chats/1/messages", nil))
	listRec := httptest.NewRecorder()
	handler.Messages(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", listRec.Code)
	}
}

func TestHandlerRejectsBadChatPath(t *testing.T) {
	handler := NewHandler(NewService(&fakeChatRepo{}))
	req := withTutor(httptest.NewRequest(http.MethodGet, "/api/chats/bad/messages", nil))
	rec := httptest.NewRecorder()

	handler.Messages(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func withTutor(req *http.Request) *http.Request {
	claims := commonauth.Claims{Role: commonauth.RoleTutor, TutorID: 1}
	return req.WithContext(commonauth.WithClaims(req.Context(), claims))
}

func withStudent(req *http.Request) *http.Request {
	claims := commonauth.Claims{Role: commonauth.RoleStudent, TutorID: 1, StudentID: 4}
	return req.WithContext(commonauth.WithClaims(req.Context(), claims))
}
