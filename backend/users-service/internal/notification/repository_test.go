package notification

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepositoryListNotifications(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	studentID := int64(4)
	lessonID := int64(7)
	mock.ExpectQuery("SELECT id, tutor_id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(notificationColumns()).
			AddRow(1, 1, studentID, lessonID, "message", "Новое сообщение", "От ученика", "tutor", false, now))
	mock.ExpectQuery("SELECT id, tutor_id").WithArgs(studentID).
		WillReturnRows(sqlmock.NewRows(notificationColumns()).
			AddRow(2, 1, studentID, nil, "cancel", "Отмена", "Занятие отменено", "student", false, now))

	repo := NewPostgresRepository(db)
	tutorItems, err := repo.ListForTutor(context.Background(), 1)
	if err != nil {
		t.Fatalf("list for tutor: %v", err)
	}
	if len(tutorItems) != 1 || tutorItems[0].StudentID == nil || *tutorItems[0].StudentID != studentID {
		t.Fatalf("unexpected tutor notifications: %+v", tutorItems)
	}

	studentItems, err := repo.ListForStudent(context.Background(), studentID)
	if err != nil {
		t.Fatalf("list for student: %v", err)
	}
	if len(studentItems) != 1 || studentItems[0].RecipientType != "student" {
		t.Fatalf("unexpected student notifications: %+v", studentItems)
	}
}

func TestRepositoryMarkRead(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectQuery("UPDATE notifications").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(notificationColumns()).
			AddRow(1, 1, nil, nil, "message", "Новое сообщение", "От ученика", "tutor", true, now))

	repo := NewPostgresRepository(db)
	item, err := repo.MarkRead(context.Background(), 1)
	if err != nil {
		t.Fatalf("mark read: %v", err)
	}
	if !item.IsRead {
		t.Fatalf("notification must be read")
	}
}

func TestRepositoryMarkReadNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("UPDATE notifications").WithArgs(int64(99)).
		WillReturnRows(sqlmock.NewRows(notificationColumns()))

	repo := NewPostgresRepository(db)
	_, err = repo.MarkRead(context.Background(), 99)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func notificationColumns() []string {
	return []string{
		"id",
		"tutor_id",
		"student_id",
		"lesson_id",
		"type",
		"title",
		"description",
		"recipient_type",
		"is_read",
		"created_at",
	}
}
