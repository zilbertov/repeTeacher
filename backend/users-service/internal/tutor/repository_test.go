package tutor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepositoryListAndGetTutors(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectQuery("SELECT t.id").
		WillReturnRows(sqlmock.NewRows(tutorColumns()).
			AddRow(1, "Вадим", "v4bem@ya.ru", "89198318673", "Заметка", now, "Математика,Русский язык"))
	mock.ExpectQuery("SELECT t.id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(tutorColumns()).
			AddRow(1, "Вадим", "v4bem@ya.ru", "89198318673", "Заметка", now, "Математика"))

	repo := NewPostgresRepository(db)
	items, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("list tutors: %v", err)
	}
	if len(items) != 1 || len(items[0].Subjects) != 2 {
		t.Fatalf("unexpected tutors: %+v", items)
	}

	item, err := repo.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("get tutor: %v", err)
	}
	if item.ID != 1 || item.Subjects[0] != "Математика" {
		t.Fatalf("unexpected tutor: %+v", item)
	}
}

func TestRepositoryGetTutorNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT t.id").WithArgs(int64(99)).
		WillReturnRows(sqlmock.NewRows(tutorColumns()))

	repo := NewPostgresRepository(db)
	_, err = repo.Get(context.Background(), 99)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRepositoryCreateTutor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO tutors").
		WithArgs("Анна", "anna@example.com", "89190000000", "demo").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	mock.ExpectExec("INSERT INTO notification_settings").WithArgs(int64(2)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO tutor_subjects").WithArgs(int64(2), "Математика").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT t.id").WithArgs(int64(2)).
		WillReturnRows(sqlmock.NewRows(tutorColumns()).
			AddRow(2, "Анна", "anna@example.com", "89190000000", "demo", now, "Математика"))

	repo := NewPostgresRepository(db)
	item, err := repo.Create(context.Background(), CreateTutorRequest{
		Name:     "Анна",
		Email:    "anna@example.com",
		Phone:    "89190000000",
		Subjects: []string{"Математика"},
		Notes:    "demo",
	})
	if err != nil {
		t.Fatalf("create tutor: %v", err)
	}
	if item.ID != 2 {
		t.Fatalf("unexpected tutor id: %d", item.ID)
	}
}

func TestRepositoryUpdateTutorNotes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectExec("UPDATE tutors").
		WithArgs("Новая заметка", int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT t.id").WithArgs(int64(2)).
		WillReturnRows(sqlmock.NewRows(tutorColumns()).
			AddRow(2, "Анна", "anna@example.com", "89190000000", "Новая заметка", now, "Математика"))

	repo := NewPostgresRepository(db)
	item, err := repo.UpdateNotes(context.Background(), 2, "Новая заметка")
	if err != nil {
		t.Fatalf("update notes: %v", err)
	}
	if item.Notes != "Новая заметка" {
		t.Fatalf("notes were not saved")
	}
}

func tutorColumns() []string {
	return []string{"id", "name", "email", "phone", "notes", "created_at", "subjects"}
}
