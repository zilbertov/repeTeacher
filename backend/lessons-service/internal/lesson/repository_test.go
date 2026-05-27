package lesson

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepositoryListLessons(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectQuery("SELECT l.id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(lessonColumns()).
			AddRow(1, 1, 2, "Вадим", "Полина", "Математика", "ОГЭ", "2026-04-02", "10:00:00", 60, "offline", true, 700, "planned", now, now))
	mock.ExpectQuery("SELECT id, lesson_id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(fileColumns()).
			AddRow(1, 1, "material", "example.JPG", "/files/example.JPG"))

	repo := NewPostgresRepository(db)
	items, err := repo.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 1 || len(items[0].Files) != 1 {
		t.Fatalf("unexpected lessons: %+v", items)
	}
}

func TestRepositoryGetNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT l.id").WithArgs(int64(99)).WillReturnRows(sqlmock.NewRows(lessonColumns()))

	repo := NewPostgresRepository(db)
	_, err = repo.Get(context.Background(), 99)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRepositoryCreateAndCancel(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO lessons").
		WithArgs(int64(1), int64(2), "Математика", "ОГЭ", "2026-04-02", "10:00", 60, "offline", true, 700).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec("INSERT INTO notifications").
		WithArgs(int64(1), int64(2), int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT l.id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(lessonColumns()).
			AddRow(1, 1, 2, "Вадим", "Полина", "Математика", "ОГЭ", "2026-04-02", "10:00:00", 60, "offline", true, 700, "planned", now, now))
	mock.ExpectQuery("SELECT id, lesson_id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(fileColumns()))
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE lessons").
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO notifications").
		WithArgs(int64(1), "Занятие отменено репетитором.", "student").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT l.id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(lessonColumns()).
			AddRow(1, 1, 2, "Вадим", "Полина", "Математика", "ОГЭ", "2026-04-02", "10:00:00", 60, "offline", true, 700, "cancelled", now, now))
	mock.ExpectQuery("SELECT id, lesson_id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(fileColumns()))

	repo := NewPostgresRepository(db)
	item, err := repo.Create(context.Background(), 1, CreateLessonRequest{
		StudentID:       2,
		Subject:         "Математика",
		ExamType:        "ОГЭ",
		LessonDate:      "2026-04-02",
		StartTime:       "10:00",
		DurationMinutes: 60,
		Format:          "offline",
		HasHomework:     true,
		Price:           700,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if item.ID != 1 {
		t.Fatalf("unexpected id")
	}

	item, err = repo.Cancel(context.Background(), 1, CancelRequest{SenderType: "tutor"})
	if err != nil {
		t.Fatalf("cancel: %v", err)
	}
	if item.Status != "cancelled" {
		t.Fatalf("lesson was not cancelled")
	}
}

func TestRepositoryAddFile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("INSERT INTO lesson_files").
		WithArgs(int64(1), "material", "file.pdf", "/files/file.pdf").
		WillReturnRows(sqlmock.NewRows(fileColumns()).AddRow(1, 1, "material", "file.pdf", "/files/file.pdf"))

	repo := NewPostgresRepository(db)
	file, err := repo.AddFile(context.Background(), 1, AddFileRequest{
		FileType: "material",
		FileName: "file.pdf",
		FilePath: "/files/file.pdf",
	})
	if err != nil {
		t.Fatalf("add file: %v", err)
	}
	if file.FileName != "file.pdf" {
		t.Fatalf("unexpected file")
	}
}

func TestRepositoryListByStudentUpdateAndReschedule(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectQuery("SELECT l.id").WithArgs(int64(2)).
		WillReturnRows(sqlmock.NewRows(lessonColumns()).
			AddRow(1, 1, 2, "Вадим", "Полина", "Математика", "ОГЭ", "2026-04-02", "10:00:00", 60, "offline", true, 700, "planned", now, now))
	mock.ExpectQuery("SELECT id, lesson_id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(fileColumns()))
	mock.ExpectExec("UPDATE lessons").
		WithArgs(int64(2), "Русский язык", "ЕГЭ", "2026-04-03", "11:00", 90, "online", false, 900, int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT l.id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(lessonColumns()).
			AddRow(1, 1, 2, "Вадим", "Полина", "Русский язык", "ЕГЭ", "2026-04-03", "11:00:00", 90, "online", false, 900, "planned", now, now))
	mock.ExpectQuery("SELECT id, lesson_id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(fileColumns()))
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE lessons").
		WithArgs("2026-04-04", "12:00", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO notifications").
		WithArgs(int64(1), "Ученик запросил перенос занятия.", "tutor").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT l.id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(lessonColumns()).
			AddRow(1, 1, 2, "Вадим", "Полина", "Русский язык", "ЕГЭ", "2026-04-04", "12:00:00", 90, "online", false, 900, "planned", now, now))
	mock.ExpectQuery("SELECT id, lesson_id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(fileColumns()))

	repo := NewPostgresRepository(db)
	items, err := repo.ListByStudent(context.Background(), 2)
	if err != nil {
		t.Fatalf("list by student: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("unexpected lessons: %+v", items)
	}

	item, err := repo.Update(context.Background(), 1, UpdateLessonRequest{
		StudentID:       2,
		Subject:         "Русский язык",
		ExamType:        "ЕГЭ",
		LessonDate:      "2026-04-03",
		StartTime:       "11:00",
		DurationMinutes: 90,
		Format:          "online",
		HasHomework:     false,
		Price:           900,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if item.Subject != "Русский язык" {
		t.Fatalf("lesson was not updated")
	}

	item, err = repo.Reschedule(context.Background(), 1, RescheduleRequest{
		LessonDate: "2026-04-04",
		StartTime:  "12:00",
		SenderType: "student",
	})
	if err != nil {
		t.Fatalf("reschedule: %v", err)
	}
	if item.LessonDate != "2026-04-04" {
		t.Fatalf("lesson was not rescheduled")
	}
}

func lessonColumns() []string {
	return []string{
		"id",
		"tutor_id",
		"student_id",
		"tutor_name",
		"student_name",
		"subject",
		"exam_type",
		"lesson_date",
		"start_time",
		"duration_minutes",
		"format",
		"has_homework",
		"price",
		"status",
		"created_at",
		"updated_at",
	}
}

func fileColumns() []string {
	return []string{"id", "lesson_id", "file_type", "file_name", "file_path"}
}
