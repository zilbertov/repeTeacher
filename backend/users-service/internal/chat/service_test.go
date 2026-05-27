package chat

import (
	"context"
	"errors"
	"testing"
)

type fakeChatRepo struct {
	messages []Message
}

func (r *fakeChatRepo) List(ctx context.Context, tutorID int64) ([]Chat, error) {
	return []Chat{{ID: 1, TutorID: tutorID, StudentID: 1, ParticipantName: "Илья Антонов"}}, nil
}

func (r *fakeChatRepo) ListByStudent(ctx context.Context, studentID int64) ([]Chat, error) {
	return []Chat{{ID: 1, TutorID: 1, StudentID: studentID, ParticipantName: "Вадим Зильбертов"}}, nil
}

func (r *fakeChatRepo) Create(ctx context.Context, tutorID int64, studentID int64) (Chat, error) {
	return Chat{ID: 1, TutorID: tutorID, StudentID: studentID, ParticipantName: "Илья Антонов"}, nil
}

func (r *fakeChatRepo) ListMessages(ctx context.Context, chatID int64) ([]Message, error) {
	return r.messages, nil
}

func (r *fakeChatRepo) SendMessage(ctx context.Context, chatID int64, senderType string, text string) (Message, error) {
	item := Message{ID: int64(len(r.messages) + 1), ChatID: chatID, SenderType: senderType, Text: text}
	r.messages = append(r.messages, item)
	return item, nil
}

func TestSendMessage(t *testing.T) {
	repo := &fakeChatRepo{}
	service := NewService(repo)

	item, err := service.SendMessage(context.Background(), 1, SendMessageRequest{Text: "Привет"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.SenderType != "tutor" {
		t.Fatalf("expected tutor sender")
	}
}

func TestSendStudentMessage(t *testing.T) {
	repo := &fakeChatRepo{}
	service := NewService(repo)

	item, err := service.SendMessage(context.Background(), 1, SendMessageRequest{Text: "Здравствуйте", SenderType: "student"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.SenderType != "student" {
		t.Fatalf("expected student sender")
	}
}

func TestCreateChat(t *testing.T) {
	service := NewService(&fakeChatRepo{})

	item, err := service.Create(context.Background(), 1, CreateChatRequest{StudentID: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.StudentID != 2 {
		t.Fatalf("wrong student id")
	}
}

func TestCreateChatCanUseTutorFromRequest(t *testing.T) {
	service := NewService(&fakeChatRepo{})

	item, err := service.Create(context.Background(), 1, CreateChatRequest{StudentID: 2, TutorID: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.TutorID != 3 {
		t.Fatalf("wrong tutor id")
	}
}

func TestListByStudentValidation(t *testing.T) {
	service := NewService(&fakeChatRepo{})

	_, err := service.ListByStudent(context.Background(), 0)
	if !errors.Is(err, ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestSendEmptyMessageReturnsError(t *testing.T) {
	service := NewService(&fakeChatRepo{})

	_, err := service.SendMessage(context.Background(), 1, SendMessageRequest{Text: "   "})
	if !errors.Is(err, ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestSendMessageRejectsUnknownSender(t *testing.T) {
	service := NewService(&fakeChatRepo{})

	_, err := service.SendMessage(context.Background(), 1, SendMessageRequest{Text: "Привет", SenderType: "admin"})
	if !errors.Is(err, ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}
