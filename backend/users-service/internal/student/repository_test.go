package student

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepositoryList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows(studentColumns()).
		AddRow(1, 1, "Илья", "ilya@email.ru", "123", "ОГЭ", "request", "", now, now, "Математика")
	mock.ExpectQuery("SELECT s.id").WithArgs(int64(1)).WillReturnRows(rows)

	repo := NewPostgresRepository(db)
	items, err := repo.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 1 || items[0].Subjects[0] != "Математика" {
		t.Fatalf("unexpected students: %+v", items)
	}
}

func TestRepositoryGetNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT s.id").WithArgs(int64(99)).WillReturnRows(sqlmock.NewRows(studentColumns()))

	repo := NewPostgresRepository(db)
	_, err = repo.Get(context.Background(), 99)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRepositoryCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO students").
		WithArgs(int64(1), "Новый", "new@email.ru", "123", "ЕГЭ", "request", "").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))
	mock.ExpectExec("INSERT INTO student_subjects").
		WithArgs(int64(10), "Математика").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT s.id").WithArgs(int64(10)).
		WillReturnRows(sqlmock.NewRows(studentColumns()).
			AddRow(10, 1, "Новый", "new@email.ru", "123", "ЕГЭ", "request", "", now, now, "Математика"))

	repo := NewPostgresRepository(db)
	item, err := repo.Create(context.Background(), 1, CreateStudentRequest{
		Name:     "Новый",
		Email:    "new@email.ru",
		Phone:    "123",
		Subjects: []string{"Математика"},
		ExamType: "ЕГЭ",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if item.ID != 10 {
		t.Fatalf("unexpected id: %d", item.ID)
	}
}

func TestRepositorySetStatusAndUpdateNotes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectExec("UPDATE students").
		WithArgs("active", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT s.id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(studentColumns()).
			AddRow(1, 1, "Илья", "ilya@email.ru", "123", "ОГЭ", "active", "", now, now, "Математика"))
	mock.ExpectExec("UPDATE students").
		WithArgs("Новая заметка", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT s.id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(studentColumns()).
			AddRow(1, 1, "Илья", "ilya@email.ru", "123", "ОГЭ", "active", "Новая заметка", now, now, "Математика"))

	repo := NewPostgresRepository(db)
	item, err := repo.SetStatus(context.Background(), 1, "active")
	if err != nil {
		t.Fatalf("set status: %v", err)
	}
	if item.Status != "active" {
		t.Fatalf("status was not changed")
	}

	item, err = repo.UpdateNotes(context.Background(), 1, "Новая заметка")
	if err != nil {
		t.Fatalf("update notes: %v", err)
	}
	if item.Notes != "Новая заметка" {
		t.Fatalf("notes were not changed")
	}
}

func TestRepositoryDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM students").WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(0, 1))

	repo := NewPostgresRepository(db)
	if err := repo.Delete(context.Background(), 1); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func studentColumns() []string {
	return []string{
		"id",
		"tutor_id",
		"name",
		"email",
		"phone",
		"exam_type",
		"status",
		"notes",
		"created_at",
		"updated_at",
		"subjects",
	}
}
