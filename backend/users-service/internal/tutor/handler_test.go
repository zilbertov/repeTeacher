package tutor

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerListAndCreateTutors(t *testing.T) {
	handler := NewHandler(NewService(&fakeTutorRepo{}))

	listReq := httptest.NewRequest(http.MethodGet, "/api/tutors", nil)
	listRec := httptest.NewRecorder()
	handler.ListOrCreate(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", listRec.Code)
	}

	body := bytes.NewBufferString(`{"name":"Анна Смирнова","email":"anna@example.com","phone":"89191234567","subjects":["Математика"]}`)
	createReq := httptest.NewRequest(http.MethodPost, "/api/tutors", body)
	createRec := httptest.NewRecorder()
	handler.ListOrCreate(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", createRec.Code, createRec.Body.String())
	}
}

func TestHandlerGetTutor(t *testing.T) {
	repo := &fakeTutorRepo{items: []Tutor{{ID: 1, Name: "Анна", Email: "anna@example.com", Subjects: []string{"Математика"}}}}
	handler := NewHandler(NewService(repo))
	req := httptest.NewRequest(http.MethodGet, "/api/tutors/1", nil)
	rec := httptest.NewRecorder()

	handler.ByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandlerUpdateTutorNotes(t *testing.T) {
	repo := &fakeTutorRepo{items: []Tutor{{ID: 1, Name: "Анна", Email: "anna@example.com", Subjects: []string{"Математика"}}}}
	handler := NewHandler(NewService(repo))
	req := httptest.NewRequest(http.MethodPost, "/api/tutors/1/notes", bytes.NewBufferString(`{"notes":"Заметка"}`))
	rec := httptest.NewRecorder()

	handler.ByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandlerTutorBadPath(t *testing.T) {
	handler := NewHandler(NewService(&fakeTutorRepo{}))
	req := httptest.NewRequest(http.MethodGet, "/api/tutors/bad", nil)
	rec := httptest.NewRecorder()

	handler.ByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
