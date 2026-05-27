package chat

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

func (s *Service) List(ctx context.Context, tutorID int64) ([]Chat, error) {
	return s.repo.List(ctx, tutorID)
}

func (s *Service) ListByStudent(ctx context.Context, studentID int64) ([]Chat, error) {
	if studentID == 0 {
		return nil, ErrBadRequest
	}
	return s.repo.ListByStudent(ctx, studentID)
}

func (s *Service) Create(ctx context.Context, tutorID int64, req CreateChatRequest) (Chat, error) {
	if req.StudentID == 0 {
		return Chat{}, ErrBadRequest
	}
	if req.TutorID != 0 {
		tutorID = req.TutorID
	}
	return s.repo.Create(ctx, tutorID, req.StudentID)
}

func (s *Service) ListMessages(ctx context.Context, chatID int64) ([]Message, error) {
	return s.repo.ListMessages(ctx, chatID)
}

func (s *Service) SendMessage(ctx context.Context, chatID int64, req SendMessageRequest) (Message, error) {
	text := strings.TrimSpace(req.Text)
	if text == "" {
		return Message{}, ErrBadRequest
	}

	senderType := strings.TrimSpace(req.SenderType)
	if senderType == "" {
		senderType = "tutor"
	}
	if senderType != "tutor" && senderType != "student" {
		return Message{}, ErrBadRequest
	}

	return s.repo.SendMessage(ctx, chatID, senderType, text)
}
