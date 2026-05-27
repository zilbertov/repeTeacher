package chat

import (
	"context"
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("chat not found")

type Repository interface {
	List(ctx context.Context, tutorID int64) ([]Chat, error)
	ListByStudent(ctx context.Context, studentID int64) ([]Chat, error)
	Create(ctx context.Context, tutorID int64, studentID int64) (Chat, error)
	ListMessages(ctx context.Context, chatID int64) ([]Message, error)
	SendMessage(ctx context.Context, chatID int64, senderType string, text string) (Message, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) List(ctx context.Context, tutorID int64) ([]Chat, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.id, c.tutor_id, c.student_id, s.name,
		       COALESCE(m.text, '') AS last_message,
		       COALESCE(m.created_at, c.created_at) AS last_message_time
		FROM chats c
		JOIN students s ON s.id = c.student_id
		LEFT JOIN LATERAL (
		    SELECT text, created_at
		    FROM messages
		    WHERE chat_id = c.id
		    ORDER BY created_at DESC
		    LIMIT 1
		) m ON true
		WHERE c.tutor_id = $1
		ORDER BY last_message_time DESC
	`, tutorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := make([]Chat, 0)
	for rows.Next() {
		var item Chat
		err := rows.Scan(&item.ID, &item.TutorID, &item.StudentID, &item.ParticipantName, &item.LastMessage, &item.LastMessageTime)
		if err != nil {
			return nil, err
		}
		chats = append(chats, item)
	}
	return chats, rows.Err()
}

func (r *PostgresRepository) ListByStudent(ctx context.Context, studentID int64) ([]Chat, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.id, c.tutor_id, c.student_id, t.name,
		       COALESCE(m.text, '') AS last_message,
		       COALESCE(m.created_at, c.created_at) AS last_message_time
		FROM chats c
		JOIN tutors t ON t.id = c.tutor_id
		LEFT JOIN LATERAL (
		    SELECT text, created_at
		    FROM messages
		    WHERE chat_id = c.id
		    ORDER BY created_at DESC
		    LIMIT 1
		) m ON true
		WHERE c.student_id = $1
		ORDER BY last_message_time DESC
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := make([]Chat, 0)
	for rows.Next() {
		var item Chat
		err := rows.Scan(&item.ID, &item.TutorID, &item.StudentID, &item.ParticipantName, &item.LastMessage, &item.LastMessageTime)
		if err != nil {
			return nil, err
		}
		chats = append(chats, item)
	}
	return chats, rows.Err()
}

func (r *PostgresRepository) Create(ctx context.Context, tutorID int64, studentID int64) (Chat, error) {
	item, err := r.getByStudent(ctx, tutorID, studentID)
	if err == nil {
		return item, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return Chat{}, err
	}

	var id int64
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO chats (tutor_id, student_id)
		VALUES ($1, $2)
		RETURNING id
	`, tutorID, studentID).Scan(&id)
	if err != nil {
		return Chat{}, err
	}

	row := r.db.QueryRowContext(ctx, `
		SELECT c.id, c.tutor_id, c.student_id, s.name, '' AS last_message, c.created_at AS last_message_time
		FROM chats c
		JOIN students s ON s.id = c.student_id
		WHERE c.id = $1
	`, id)

	var created Chat
	err = row.Scan(&created.ID, &created.TutorID, &created.StudentID, &created.ParticipantName, &created.LastMessage, &created.LastMessageTime)
	return created, err
}

func (r *PostgresRepository) getByStudent(ctx context.Context, tutorID int64, studentID int64) (Chat, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT c.id, c.tutor_id, c.student_id, s.name,
		       COALESCE(m.text, '') AS last_message,
		       COALESCE(m.created_at, c.created_at) AS last_message_time
		FROM chats c
		JOIN students s ON s.id = c.student_id
		LEFT JOIN LATERAL (
		    SELECT text, created_at
		    FROM messages
		    WHERE chat_id = c.id
		    ORDER BY created_at DESC
		    LIMIT 1
		) m ON true
		WHERE c.tutor_id = $1 AND c.student_id = $2
		ORDER BY c.id
		LIMIT 1
	`, tutorID, studentID)

	var item Chat
	err := row.Scan(&item.ID, &item.TutorID, &item.StudentID, &item.ParticipantName, &item.LastMessage, &item.LastMessageTime)
	return item, err
}

func (r *PostgresRepository) ListMessages(ctx context.Context, chatID int64) ([]Message, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, chat_id, sender_type, text, created_at
		FROM messages
		WHERE chat_id = $1
		ORDER BY created_at
	`, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]Message, 0)
	for rows.Next() {
		var item Message
		err := rows.Scan(&item.ID, &item.ChatID, &item.SenderType, &item.Text, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, item)
	}
	return messages, rows.Err()
}

func (r *PostgresRepository) SendMessage(ctx context.Context, chatID int64, senderType string, text string) (Message, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Message{}, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `
		INSERT INTO messages (chat_id, sender_type, text)
		VALUES ($1, $2, $3)
		RETURNING id, chat_id, sender_type, text, created_at
	`, chatID, senderType, text)

	var item Message
	err = row.Scan(&item.ID, &item.ChatID, &item.SenderType, &item.Text, &item.CreatedAt)
	if err != nil {
		return Message{}, err
	}

	recipientType := "student"
	description := "Новое сообщение от репетитора."
	if senderType == "student" {
		recipientType = "tutor"
		description = "Новое сообщение от ученика."
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO notifications (tutor_id, student_id, type, title, description, recipient_type)
		SELECT tutor_id, student_id, 'message', 'Новое сообщение', $2, $3
		FROM chats
		WHERE id = $1
	`, chatID, description, recipientType)
	if err != nil {
		return Message{}, err
	}

	if err := tx.Commit(); err != nil {
		return Message{}, err
	}
	return item, nil
}
