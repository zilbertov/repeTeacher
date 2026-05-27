package lesson

import (
	"context"
	"errors"
	"testing"
)

type fakeLessonRepo struct {
	items map[int64]Lesson
}

func newFakeLessonRepo() *fakeLessonRepo {
	return &fakeLessonRepo{items: map[int64]Lesson{
		1: {ID: 1, TutorID: 1, StudentID: 2, TutorName: "Вадим", StudentName: "Полина", Subject: "Математика", LessonDate: "2026-04-02", StartTime: "10:00:00", Format: "offline", DurationMinutes: 60, Price: 700, Status: "planned"},
	}}
}

func (r *fakeLessonRepo) List(ctx context.Context, tutorID int64) ([]Lesson, error) {
	return []Lesson{r.items[1]}, nil
}

func (r *fakeLessonRepo) ListByStudent(ctx context.Context, studentID int64) ([]Lesson, error) {
	return []Lesson{r.items[1]}, nil
}

func (r *fakeLessonRepo) Get(ctx context.Context, id int64) (Lesson, error) {
	item, ok := r.items[id]
	if !ok {
		return Lesson{}, ErrNotFound
	}
	return item, nil
}

func (r *fakeLessonRepo) Create(ctx context.Context, tutorID int64, req CreateLessonRequest) (Lesson, error) {
	item := Lesson{ID: 2, TutorID: tutorID, StudentID: req.StudentID, Subject: req.Subject, LessonDate: req.LessonDate, StartTime: req.StartTime, Format: req.Format, DurationMinutes: req.DurationMinutes, Price: req.Price}
	r.items[item.ID] = item
	return item, nil
}

func (r *fakeLessonRepo) Update(ctx context.Context, id int64, req UpdateLessonRequest) (Lesson, error) {
	item := r.items[id]
	item.Subject = req.Subject
	r.items[id] = item
	return item, nil
}

func (r *fakeLessonRepo) Reschedule(ctx context.Context, id int64, req RescheduleRequest) (Lesson, error) {
	item := r.items[id]
	item.LessonDate = req.LessonDate
	item.StartTime = req.StartTime
	r.items[id] = item
	return item, nil
}

func (r *fakeLessonRepo) Cancel(ctx context.Context, id int64, req CancelRequest) (Lesson, error) {
	item := r.items[id]
	item.Status = "cancelled"
	r.items[id] = item
	return item, nil
}

func (r *fakeLessonRepo) AddFile(ctx context.Context, id int64, req AddFileRequest) (LessonFile, error) {
	return LessonFile{ID: 1, LessonID: id, FileType: req.FileType, FileName: req.FileName, FilePath: req.FilePath}, nil
}

func TestCreateLessonValidation(t *testing.T) {
	service := NewService(newFakeLessonRepo())

	_, err := service.Create(context.Background(), 1, CreateLessonRequest{})
	if !errors.Is(err, ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestCancelLesson(t *testing.T) {
	service := NewService(newFakeLessonRepo())

	item, err := service.Cancel(context.Background(), 1, CancelRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Status != "cancelled" {
		t.Fatalf("expected cancelled lesson")
	}
}

func TestCancelLessonRejectsUnknownSender(t *testing.T) {
	service := NewService(newFakeLessonRepo())

	_, err := service.Cancel(context.Background(), 1, CancelRequest{SenderType: "admin"})
	if !errors.Is(err, ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestRescheduleLesson(t *testing.T) {
	service := NewService(newFakeLessonRepo())

	item, err := service.Reschedule(context.Background(), 1, RescheduleRequest{LessonDate: "2026-05-01", StartTime: "12:00"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.LessonDate != "2026-05-01" {
		t.Fatalf("lesson date was not changed")
	}
}

func TestListByStudentValidation(t *testing.T) {
	service := NewService(newFakeLessonRepo())

	_, err := service.ListByStudent(context.Background(), 0)
	if !errors.Is(err, ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestAddFileValidation(t *testing.T) {
	service := NewService(newFakeLessonRepo())

	_, err := service.AddFile(context.Background(), 1, AddFileRequest{FileType: "bad"})
	if !errors.Is(err, ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}
