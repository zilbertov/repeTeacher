package profile

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

func (s *Service) Get(ctx context.Context, tutorID int64) (Profile, error) {
	return s.repo.Get(ctx, tutorID)
}

func (s *Service) Update(ctx context.Context, tutorID int64, req UpdateProfileRequest) (Profile, error) {
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" {
		return Profile{}, ErrBadRequest
	}
	if req.Subjects != nil {
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
		req.Subjects = subjects
	}
	return s.repo.Update(ctx, tutorID, req)
}

func (s *Service) ChangePassword(ctx context.Context, tutorID int64, req ChangePasswordRequest) error {
	if len(req.NewPassword) < 4 {
		return ErrBadRequest
	}

	passwordHash := "changed:" + req.NewPassword
	return s.repo.ChangePassword(ctx, tutorID, passwordHash)
}

func (s *Service) GetSettings(ctx context.Context, tutorID int64) (NotificationSettings, error) {
	return s.repo.GetSettings(ctx, tutorID)
}

func (s *Service) UpdateSettings(ctx context.Context, tutorID int64, settings NotificationSettings) (NotificationSettings, error) {
	return s.repo.UpdateSettings(ctx, tutorID, settings)
}
