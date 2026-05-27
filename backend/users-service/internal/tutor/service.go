package tutor

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

func (s *Service) List(ctx context.Context) ([]Tutor, error) {
	return s.repo.List(ctx)
}

func (s *Service) Get(ctx context.Context, id int64) (Tutor, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Create(ctx context.Context, req CreateTutorRequest) (Tutor, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Notes = strings.TrimSpace(req.Notes)
	if req.Name == "" || req.Email == "" {
		return Tutor{}, ErrBadRequest
	}

	subjects := make([]string, 0, len(req.Subjects))
	seen := make(map[string]bool)
	for _, subject := range req.Subjects {
		subject = strings.TrimSpace(subject)
		if subject == "" || seen[subject] {
			continue
		}
		seen[subject] = true
		subjects = append(subjects, subject)
	}
	if len(subjects) == 0 {
		return Tutor{}, ErrBadRequest
	}
	req.Subjects = subjects

	return s.repo.Create(ctx, req)
}

func (s *Service) UpdateNotes(ctx context.Context, id int64, req UpdateTutorNotesRequest) (Tutor, error) {
	return s.repo.UpdateNotes(ctx, id, strings.TrimSpace(req.Notes))
}
