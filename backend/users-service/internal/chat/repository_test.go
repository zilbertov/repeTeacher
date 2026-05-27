package chat

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepositoryListChats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectQuery("SELECT c.id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(chatColumns()).
			AddRow(1, 1, 2, "Полина", "Привет", now))

	repo := NewPostgresRepository(db)
	items, err := repo.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 1 || items[0].ParticipantName != "Полина" {
		t.Fatalf("unexpected chats: %+v", items)
	}
}

func TestRepositoryListChatsByStudent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectQuery("SELECT c.id").WithArgs(int64(2)).
		WillReturnRows(sqlmock.NewRows(chatColumns()).
			AddRow(1, 1, 2, "Вадим", "", now))

	repo := NewPostgresRepository(db)
	items, err := repo.ListByStudent(context.Background(), 2)
	if err != nil {
		t.Fatalf("list by student: %v", err)
	}
	if len(items) != 1 || items[0].ParticipantName != "Вадим" {
		t.Fatalf("unexpected chats: %+v", items)
	}
}

func TestRepositoryCreateChat(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectQuery("SELECT c.id").WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows(chatColumns()))
	mock.ExpectQuery("INSERT INTO chats").WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))
	mock.ExpectQuery("SELECT c.id").WithArgs(int64(10)).
		WillReturnRows(sqlmock.NewRows(chatColumns()).
			AddRow(10, 1, 2, "Полина", "", now))

	repo := NewPostgresRepository(db)
	item, err := repo.Create(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if item.ID != 10 || item.StudentID != 2 {
		t.Fatalf("unexpected chat: %+v", item)
	}
}

func TestRepositorySendMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO messages").WithArgs(int64(1), "student", "Здравствуйте").
		WillReturnRows(sqlmock.NewRows([]string{"id", "chat_id", "sender_type", "text", "created_at"}).
			AddRow(3, 1, "student", "Здравствуйте", now))
	mock.ExpectExec("INSERT INTO notifications").WithArgs(int64(1), "Новое сообщение от ученика.", "tutor").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewPostgresRepository(db)
	item, err := repo.SendMessage(context.Background(), 1, "student", "Здравствуйте")
	if err != nil {
		t.Fatalf("send message: %v", err)
	}
	if item.SenderType != "student" {
		t.Fatalf("unexpected message: %+v", item)
	}
}

func chatColumns() []string {
	return []string{"id", "tutor_id", "student_id", "participant_name", "last_message", "last_message_time"}
}
