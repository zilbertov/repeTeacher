package student

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

func (s *Service) List(ctx context.Context, tutorID int64) ([]Student, error) {
	return s.repo.List(ctx, tutorID)
}

func (s *Service) Get(ctx context.Context, id int64) (Student, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Create(ctx context.Context, tutorID int64, req CreateStudentRequest) (Student, error) {
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" {
		return Student{}, ErrBadRequest
	}
	if len(req.Subjects) == 0 {
		return Student{}, ErrBadRequest
	}
	return s.repo.Create(ctx, tutorID, req)
}

func (s *Service) Update(ctx context.Context, id int64, req UpdateStudentRequest) (Student, error) {
	if strings.TrimSpace(req.Name) == "" {
		return Student{}, ErrBadRequest
	}
	return s.repo.Update(ctx, id, req)
}

func (s *Service) Accept(ctx context.Context, id int64) (Student, error) {
	item, err := s.repo.Get(ctx, id)
	if err != nil {
		return Student{}, err
	}
	if item.Status != "request" {
		return Student{}, ErrBadRequest
	}
	return s.repo.SetStatus(ctx, id, "active")
}

func (s *Service) Archive(ctx context.Context, id int64) (Student, error) {
	return s.repo.SetStatus(ctx, id, "archived")
}

func (s *Service) UpdateNotes(ctx context.Context, id int64, notes string) (Student, error) {
	return s.repo.UpdateNotes(ctx, id, notes)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
