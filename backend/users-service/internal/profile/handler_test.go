package profile

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	commonauth "github.com/zilbertov/repe-teacher-common/auth"
)

func TestHandlerGetAndUpdateProfile(t *testing.T) {
	handler := NewHandler(NewService(newFakeProfileRepo()))

	getReq := withTutor(httptest.NewRequest(http.MethodGet, "/api/profile", nil))
	getRec := httptest.NewRecorder()
	handler.Profile(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", getRec.Code)
	}

	body := bytes.NewBufferString(`{"name":"Вадим","email":"v4bem@ya.ru","phone":"89198318673","subjects":["Математика"]}`)
	updateReq := withTutor(httptest.NewRequest(http.MethodPut, "/api/profile", body))
	updateRec := httptest.NewRecorder()
	handler.Profile(updateRec, updateReq)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", updateRec.Code, updateRec.Body.String())
	}
}

func TestHandlerChangePassword(t *testing.T) {
	handler := NewHandler(NewService(newFakeProfileRepo()))
	body := bytes.NewBufferString(`{"current_password":"old","new_password":"1234"}`)
	req := withTutor(httptest.NewRequest(http.MethodPost, "/api/profile/password", body))
	rec := httptest.NewRecorder()

	handler.Password(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandlerSettings(t *testing.T) {
	handler := NewHandler(NewService(newFakeProfileRepo()))

	getReq := withTutor(httptest.NewRequest(http.MethodGet, "/api/settings/notifications", nil))
	getRec := httptest.NewRecorder()
	handler.Settings(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", getRec.Code)
	}

	body := bytes.NewBufferString(`{"push_enabled":true,"telegram_enabled":false,"sound_enabled":true,"lesson_reminders_enabled":true}`)
	updateReq := withTutor(httptest.NewRequest(http.MethodPut, "/api/settings/notifications", body))
	updateRec := httptest.NewRecorder()
	handler.Settings(updateRec, updateReq)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", updateRec.Code, updateRec.Body.String())
	}
}

func TestHandlerProfileValidationError(t *testing.T) {
	handler := NewHandler(NewService(newFakeProfileRepo()))
	req := withTutor(httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewBufferString(`{"name":"","email":""}`)))
	rec := httptest.NewRecorder()

	handler.Profile(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func withTutor(req *http.Request) *http.Request {
	claims := commonauth.Claims{Role: commonauth.RoleTutor, TutorID: 1}
	return req.WithContext(commonauth.WithClaims(req.Context(), claims))
}
