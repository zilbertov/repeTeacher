package notification

import (
	"context"
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("notification not found")

type Repository interface {
	ListForTutor(ctx context.Context, tutorID int64) ([]Notification, error)
	ListForStudent(ctx context.Context, studentID int64) ([]Notification, error)
	MarkRead(ctx context.Context, id int64) (Notification, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ListForTutor(ctx context.Context, tutorID int64) ([]Notification, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tutor_id, student_id, lesson_id, type, title, description, recipient_type, is_read, created_at
		FROM notifications
		WHERE tutor_id = $1 AND recipient_type = 'tutor'
		ORDER BY created_at DESC, id DESC
	`, tutorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Notification, 0)
	for rows.Next() {
		item, err := scanNotification(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) ListForStudent(ctx context.Context, studentID int64) ([]Notification, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tutor_id, student_id, lesson_id, type, title, description, recipient_type, is_read, created_at
		FROM notifications
		WHERE student_id = $1 AND recipient_type = 'student'
		ORDER BY created_at DESC, id DESC
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Notification, 0)
	for rows.Next() {
		item, err := scanNotification(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) MarkRead(ctx context.Context, id int64) (Notification, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE notifications
		SET is_read = true
		WHERE id = $1
		RETURNING id, tutor_id, student_id, lesson_id, type, title, description, recipient_type, is_read, created_at
	`, id)

	item, err := scanNotification(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Notification{}, ErrNotFound
	}
	return item, err
}

type scanner interface {
	Scan(dest ...any) error
}

func scanNotification(row scanner) (Notification, error) {
	var item Notification
	var studentID sql.NullInt64
	var lessonID sql.NullInt64
	err := row.Scan(
		&item.ID,
		&item.TutorID,
		&studentID,
		&lessonID,
		&item.Type,
		&item.Title,
		&item.Description,
		&item.RecipientType,
		&item.IsRead,
		&item.CreatedAt,
	)
	if err != nil {
		return Notification{}, err
	}
	if studentID.Valid {
		item.StudentID = &studentID.Int64
	}
	if lessonID.Valid {
		item.LessonID = &lessonID.Int64
	}
	return item, nil
}
