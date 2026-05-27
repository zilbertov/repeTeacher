package session

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerLogin(t *testing.T) {
	handler := NewHandler(NewService(&fakeSessionRepo{}, "test-secret"))
	body := bytes.NewBufferString(`{"role":"tutor","email":"v4bem@ya.ru","password":"demo"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandlerLoginRejectsBadPassword(t *testing.T) {
	handler := NewHandler(NewService(&fakeSessionRepo{}, "test-secret"))
	body := bytes.NewBufferString(`{"role":"tutor","email":"v4bem@ya.ru","password":"bad"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestHandlerLoginRejectsBadMethod(t *testing.T) {
	handler := NewHandler(NewService(&fakeSessionRepo{}, "test-secret"))
	req := httptest.NewRequest(http.MethodGet, "/api/auth/login", nil)
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
