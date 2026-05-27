package notification

import (
	"net/http"
	"net/http/httptest"
	"testing"

	commonauth "github.com/zilbertov/repe-teacher-common/auth"
)

func TestHandlerListNotifications(t *testing.T) {
	handler := NewHandler(NewService(&fakeNotificationRepo{item: Notification{ID: 1, Type: "message"}}))
	req := withTutor(httptest.NewRequest(http.MethodGet, "/api/notifications", nil))
	rec := httptest.NewRecorder()

	handler.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandlerListStudentNotifications(t *testing.T) {
	handler := NewHandler(NewService(&fakeNotificationRepo{item: Notification{ID: 1, Type: "message"}}))
	req := withStudent(httptest.NewRequest(http.MethodGet, "/api/notifications?student_id=4", nil))
	rec := httptest.NewRecorder()

	handler.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandlerNotificationActions(t *testing.T) {
	handler := NewHandler(NewService(&fakeNotificationRepo{item: Notification{ID: 1}}))

	for _, action := range []string{"read", "approve", "reject"} {
		req := withTutor(httptest.NewRequest(http.MethodPost, "/api/notifications/1/"+action, nil))
		rec := httptest.NewRecorder()
		handler.ByID(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 for %s, got %d", action, rec.Code)
		}
	}
}

func TestHandlerNotificationBadPath(t *testing.T) {
	handler := NewHandler(NewService(&fakeNotificationRepo{}))
	req := withTutor(httptest.NewRequest(http.MethodPost, "/api/notifications/bad/read", nil))
	rec := httptest.NewRecorder()

	handler.ByID(rec, req)

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
