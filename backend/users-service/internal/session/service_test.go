package session

import (
	"context"
	"errors"
	"testing"
)

type fakeSessionRepo struct{}

func (r *fakeSessionRepo) FindTutorByEmail(ctx context.Context, email string) (int64, error) {
	if email != "v4bem@ya.ru" {
		return 0, ErrNotFound
	}
	return 1, nil
}

func (r *fakeSessionRepo) FindStudentByEmail(ctx context.Context, email string) (int64, int64, error) {
	if email != "student.demo@example.com" {
		return 0, 0, ErrNotFound
	}
	return 4, 1, nil
}

func TestLoginTutor(t *testing.T) {
	service := NewService(&fakeSessionRepo{}, "test-secret")

	item, err := service.Login(context.Background(), LoginRequest{Role: "tutor", Email: "v4bem@ya.ru", Password: "demo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Token == "" || item.TutorID != 1 || item.Role != "tutor" {
		t.Fatalf("unexpected login response: %+v", item)
	}
}

func TestLoginStudent(t *testing.T) {
	service := NewService(&fakeSessionRepo{}, "test-secret")

	item, err := service.Login(context.Background(), LoginRequest{Role: "student", Email: "student.demo@example.com", Password: "demo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Token == "" || item.StudentID != 4 || item.TutorID != 1 || item.Role != "student" {
		t.Fatalf("unexpected login response: %+v", item)
	}
}

func TestLoginRejectsBadPassword(t *testing.T) {
	service := NewService(&fakeSessionRepo{}, "test-secret")

	_, err := service.Login(context.Background(), LoginRequest{Role: "tutor", Email: "v4bem@ya.ru", Password: "bad"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}
