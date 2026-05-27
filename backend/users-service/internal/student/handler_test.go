package student

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	commonauth "github.com/zilbertov/repe-teacher-common/auth"
)

func TestHandlerListStudents(t *testing.T) {
	handler := NewHandler(NewService(newFakeStudentRepo()))
	req := withTutor(httptest.NewRequest(http.MethodGet, "/api/students", nil))
	rec := httptest.NewRecorder()

	handler.ListOrCreate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandlerCreateStudent(t *testing.T) {
	handler := NewHandler(NewService(newFakeStudentRepo()))
	body := bytes.NewBufferString(`{"name":"Новый","email":"new@email.ru","phone":"123","subjects":["Математика"],"exam_type":"ЕГЭ","status":"request"}`)
	req := withTutor(httptest.NewRequest(http.MethodPost, "/api/students", body))
	rec := httptest.NewRecorder()

	handler.ListOrCreate(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandlerGetAndAcceptStudent(t *testing.T) {
	handler := NewHandler(NewService(newFakeStudentRepo()))

	getReq := withTutor(httptest.NewRequest(http.MethodGet, "/api/students/1", nil))
	getRec := httptest.NewRecorder()
	handler.ByID(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", getRec.Code)
	}

	acceptReq := withTutor(httptest.NewRequest(http.MethodPost, "/api/students/1/accept", nil))
	acceptRec := httptest.NewRecorder()
	handler.ByID(acceptRec, acceptReq)
	if acceptRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", acceptRec.Code)
	}
}

func TestHandlerNotesAndDeleteStudent(t *testing.T) {
	handler := NewHandler(NewService(newFakeStudentRepo()))

	notesReq := withTutor(httptest.NewRequest(http.MethodPost, "/api/students/2/notes", bytes.NewBufferString(`{"notes":"Текст"}`)))
	notesRec := httptest.NewRecorder()
	handler.ByID(notesRec, notesReq)
	if notesRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", notesRec.Code)
	}

	deleteReq := withTutor(httptest.NewRequest(http.MethodDelete, "/api/students/2", nil))
	deleteRec := httptest.NewRecorder()
	handler.ByID(deleteRec, deleteReq)
	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", deleteRec.Code)
	}
}

func TestHandlerUpdateAndArchiveStudent(t *testing.T) {
	handler := NewHandler(NewService(newFakeStudentRepo()))

	updateBody := bytes.NewBufferString(`{"name":"Полина","email":"polina@email.ru","phone":"123","subjects":["Русский язык"],"exam_type":"ЕГЭ","status":"active","notes":""}`)
	updateReq := withTutor(httptest.NewRequest(http.MethodPut, "/api/students/2", updateBody))
	updateRec := httptest.NewRecorder()
	handler.ByID(updateRec, updateReq)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", updateRec.Code)
	}

	archiveReq := withTutor(httptest.NewRequest(http.MethodPost, "/api/students/2/archive", nil))
	archiveRec := httptest.NewRecorder()
	handler.ByID(archiveRec, archiveReq)
	if archiveRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", archiveRec.Code)
	}
}

func withTutor(req *http.Request) *http.Request {
	claims := commonauth.Claims{Role: commonauth.RoleTutor, TutorID: 1}
	return req.WithContext(commonauth.WithClaims(req.Context(), claims))
}

func TestParseStudentPath(t *testing.T) {
	id, action, ok := parseStudentPath("/api/students/12/archive")
	if !ok || id != 12 || action != "archive" {
		t.Fatalf("unexpected parse result")
	}
}
