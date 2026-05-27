package lesson

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	commonauth "github.com/zilbertov/repe-teacher-common/auth"
)

func TestHandlerListLessons(t *testing.T) {
	handler := NewHandler(NewService(newFakeLessonRepo()))
	req := withTutor(httptest.NewRequest(http.MethodGet, "/api/lessons", nil))
	rec := httptest.NewRecorder()

	handler.ListOrCreate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandlerCreateLesson(t *testing.T) {
	handler := NewHandler(NewService(newFakeLessonRepo()))
	body := bytes.NewBufferString(`{"student_id":2,"subject":"Математика","exam_type":"ОГЭ","lesson_date":"2026-04-02","start_time":"10:00","duration_minutes":60,"format":"offline","has_homework":true,"price":700}`)
	req := withTutor(httptest.NewRequest(http.MethodPost, "/api/lessons", body))
	rec := httptest.NewRecorder()

	handler.ListOrCreate(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandlerGetRescheduleCancelLesson(t *testing.T) {
	handler := NewHandler(NewService(newFakeLessonRepo()))

	getReq := withTutor(httptest.NewRequest(http.MethodGet, "/api/lessons/1", nil))
	getRec := httptest.NewRecorder()
	handler.ByID(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", getRec.Code)
	}

	rescheduleReq := withTutor(httptest.NewRequest(http.MethodPost, "/api/lessons/1/reschedule", bytes.NewBufferString(`{"lesson_date":"2026-05-01","start_time":"12:00"}`)))
	rescheduleRec := httptest.NewRecorder()
	handler.ByID(rescheduleRec, rescheduleReq)
	if rescheduleRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rescheduleRec.Code)
	}

	cancelReq := withTutor(httptest.NewRequest(http.MethodPost, "/api/lessons/1/cancel", nil))
	cancelRec := httptest.NewRecorder()
	handler.ByID(cancelRec, cancelReq)
	if cancelRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", cancelRec.Code)
	}
}

func TestHandlerAddFile(t *testing.T) {
	handler := NewHandler(NewService(newFakeLessonRepo()))
	body := bytes.NewBufferString(`{"file_type":"material","file_name":"file.pdf","file_path":"/files/file.pdf"}`)
	req := withTutor(httptest.NewRequest(http.MethodPost, "/api/lessons/1/files", body))
	rec := httptest.NewRecorder()

	handler.ByID(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func withTutor(req *http.Request) *http.Request {
	claims := commonauth.Claims{Role: commonauth.RoleTutor, TutorID: 1}
	return req.WithContext(commonauth.WithClaims(req.Context(), claims))
}

func TestParseLessonPath(t *testing.T) {
	id, action, ok := parseLessonPath("/api/lessons/7/files")
	if !ok || id != 7 || action != "files" {
		t.Fatalf("unexpected parse result")
	}
}
