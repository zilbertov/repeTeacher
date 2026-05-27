package session

import (
	"context"
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("account not found")

type Repository interface {
	FindTutorByEmail(ctx context.Context, email string) (int64, error)
	FindStudentByEmail(ctx context.Context, email string) (int64, int64, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) FindTutorByEmail(ctx context.Context, email string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `SELECT id FROM tutors WHERE email = $1`, email).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	}
	return id, err
}

func (r *PostgresRepository) FindStudentByEmail(ctx context.Context, email string) (int64, int64, error) {
	var studentID int64
	var tutorID int64
	err := r.db.QueryRowContext(ctx, `SELECT id, tutor_id FROM students WHERE email = $1`, email).Scan(&studentID, &tutorID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, 0, ErrNotFound
	}
	return studentID, tutorID, err
}
