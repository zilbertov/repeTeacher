package session

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepositoryFindTutorByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id FROM tutors").WithArgs("v4bem@ya.ru").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	repo := NewPostgresRepository(db)
	id, err := repo.FindTutorByEmail(context.Background(), "v4bem@ya.ru")
	if err != nil {
		t.Fatalf("find tutor: %v", err)
	}
	if id != 1 {
		t.Fatalf("unexpected tutor id: %d", id)
	}
}

func TestRepositoryFindStudentByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, tutor_id FROM students").WithArgs("student.demo@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "tutor_id"}).AddRow(4, 1))

	repo := NewPostgresRepository(db)
	studentID, tutorID, err := repo.FindStudentByEmail(context.Background(), "student.demo@example.com")
	if err != nil {
		t.Fatalf("find student: %v", err)
	}
	if studentID != 4 || tutorID != 1 {
		t.Fatalf("unexpected ids: student=%d tutor=%d", studentID, tutorID)
	}
}

func TestRepositoryFindTutorByEmailNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id FROM tutors").WithArgs("missing@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	repo := NewPostgresRepository(db)
	_, err = repo.FindTutorByEmail(context.Background(), "missing@example.com")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
