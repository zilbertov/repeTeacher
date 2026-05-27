package lesson

import (
	"context"
	"errors"
	"strings"
)

var ErrBadRequest = errors.New("bad request")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, tutorID int64) ([]Lesson, error) {
	return s.repo.List(ctx, tutorID)
}

func (s *Service) ListByStudent(ctx context.Context, studentID int64) ([]Lesson, error) {
	if studentID == 0 {
		return nil, ErrBadRequest
	}
	return s.repo.ListByStudent(ctx, studentID)
}

func (s *Service) Get(ctx context.Context, id int64) (Lesson, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Create(ctx context.Context, tutorID int64, req CreateLessonRequest) (Lesson, error) {
	if err := validateLesson(req); err != nil {
		return Lesson{}, err
	}
	return s.repo.Create(ctx, tutorID, req)
}

func (s *Service) Update(ctx context.Context, id int64, req UpdateLessonRequest) (Lesson, error) {
	if err := validateLesson(req); err != nil {
		return Lesson{}, err
	}
	return s.repo.Update(ctx, id, req)
}

func (s *Service) Reschedule(ctx context.Context, id int64, req RescheduleRequest) (Lesson, error) {
	if strings.TrimSpace(req.LessonDate) == "" || strings.TrimSpace(req.StartTime) == "" {
		return Lesson{}, ErrBadRequest
	}
	senderType, err := normalizeSenderType(req.SenderType)
	if err != nil {
		return Lesson{}, err
	}
	req.SenderType = senderType
	return s.repo.Reschedule(ctx, id, req)
}

func (s *Service) Cancel(ctx context.Context, id int64, req CancelRequest) (Lesson, error) {
	senderType, err := normalizeSenderType(req.SenderType)
	if err != nil {
		return Lesson{}, err
	}
	req.SenderType = senderType
	return s.repo.Cancel(ctx, id, req)
}

func (s *Service) AddFile(ctx context.Context, id int64, req AddFileRequest) (LessonFile, error) {
	if strings.TrimSpace(req.FileName) == "" || strings.TrimSpace(req.FilePath) == "" {
		return LessonFile{}, ErrBadRequest
	}
	if req.FileType != "material" && req.FileType != "homework" {
		return LessonFile{}, ErrBadRequest
	}
	return s.repo.AddFile(ctx, id, req)
}

func normalizeSenderType(senderType string) (string, error) {
	senderType = strings.TrimSpace(senderType)
	if senderType == "" {
		return "tutor", nil
	}
	if senderType != "tutor" && senderType != "student" {
		return "", ErrBadRequest
	}
	return senderType, nil
}

func validateLesson(req CreateLessonRequest) error {
	if req.StudentID == 0 {
		return ErrBadRequest
	}
	if strings.TrimSpace(req.Subject) == "" || strings.TrimSpace(req.LessonDate) == "" || strings.TrimSpace(req.StartTime) == "" {
		return ErrBadRequest
	}
	if req.Format != "online" && req.Format != "offline" {
		return ErrBadRequest
	}
	if req.DurationMinutes <= 0 || req.Price < 0 {
		return ErrBadRequest
	}
	return nil
}
