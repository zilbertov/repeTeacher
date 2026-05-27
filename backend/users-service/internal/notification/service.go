package notification

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, tutorID int64) ([]Notification, error) {
	return s.repo.ListForTutor(ctx, tutorID)
}

func (s *Service) ListForStudent(ctx context.Context, studentID int64) ([]Notification, error) {
	return s.repo.ListForStudent(ctx, studentID)
}

func (s *Service) MarkRead(ctx context.Context, id int64) (Notification, error) {
	return s.repo.MarkRead(ctx, id)
}

func (s *Service) Approve(ctx context.Context, id int64) (Notification, error) {
	return s.repo.MarkRead(ctx, id)
}

func (s *Service) Reject(ctx context.Context, id int64) (Notification, error) {
	return s.repo.MarkRead(ctx, id)
}
