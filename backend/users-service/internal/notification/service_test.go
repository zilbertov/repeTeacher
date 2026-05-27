package notification

import (
	"context"
	"testing"
)

type fakeNotificationRepo struct {
	item Notification
}

func (r *fakeNotificationRepo) ListForTutor(ctx context.Context, tutorID int64) ([]Notification, error) {
	return []Notification{r.item}, nil
}

func (r *fakeNotificationRepo) ListForStudent(ctx context.Context, studentID int64) ([]Notification, error) {
	return []Notification{r.item}, nil
}

func (r *fakeNotificationRepo) MarkRead(ctx context.Context, id int64) (Notification, error) {
	r.item.IsRead = true
	return r.item, nil
}

func TestMarkRead(t *testing.T) {
	service := NewService(&fakeNotificationRepo{item: Notification{ID: 1, IsRead: false}})

	item, err := service.MarkRead(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !item.IsRead {
		t.Fatalf("notification must be read")
	}
}
