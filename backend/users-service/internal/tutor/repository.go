package tutor

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

var ErrNotFound = errors.New("tutor not found")

type Repository interface {
	List(ctx context.Context) ([]Tutor, error)
	Get(ctx context.Context, id int64) (Tutor, error)
	Create(ctx context.Context, req CreateTutorRequest) (Tutor, error)
	UpdateNotes(ctx context.Context, id int64, notes string) (Tutor, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) List(ctx context.Context) ([]Tutor, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT t.id, t.name, t.email, t.phone, t.notes, t.created_at,
		       COALESCE(string_agg(ts.subject, ',' ORDER BY ts.id), '')
		FROM tutors t
		LEFT JOIN tutor_subjects ts ON ts.tutor_id = t.id
		GROUP BY t.id
		ORDER BY t.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Tutor, 0)
	for rows.Next() {
		item, err := scanTutor(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) Get(ctx context.Context, id int64) (Tutor, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT t.id, t.name, t.email, t.phone, t.notes, t.created_at,
		       COALESCE(string_agg(ts.subject, ',' ORDER BY ts.id), '')
		FROM tutors t
		LEFT JOIN tutor_subjects ts ON ts.tutor_id = t.id
		WHERE t.id = $1
		GROUP BY t.id
	`, id)

	item, err := scanTutor(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Tutor{}, ErrNotFound
	}
	return item, err
}

func (r *PostgresRepository) Create(ctx context.Context, req CreateTutorRequest) (Tutor, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Tutor{}, err
	}
	defer tx.Rollback()

	var id int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO tutors (name, email, phone, password_hash, notes)
		VALUES ($1, $2, $3, 'demo-password-hash', $4)
		RETURNING id
	`, req.Name, req.Email, req.Phone, req.Notes).Scan(&id)
	if err != nil {
		return Tutor{}, err
	}

	if _, err = tx.ExecContext(ctx, `INSERT INTO notification_settings (tutor_id) VALUES ($1)`, id); err != nil {
		return Tutor{}, err
	}

	for _, subject := range req.Subjects {
		subject = strings.TrimSpace(subject)
		if subject == "" {
			continue
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO tutor_subjects (tutor_id, subject) VALUES ($1, $2)`, id, subject); err != nil {
			return Tutor{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return Tutor{}, err
	}
	return r.Get(ctx, id)
}

func (r *PostgresRepository) UpdateNotes(ctx context.Context, id int64, notes string) (Tutor, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE tutors
		SET notes = $1
		WHERE id = $2
	`, notes, id)
	if err != nil {
		return Tutor{}, err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return Tutor{}, ErrNotFound
	}
	return r.Get(ctx, id)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanTutor(row scanner) (Tutor, error) {
	var item Tutor
	var subjects string
	err := row.Scan(&item.ID, &item.Name, &item.Email, &item.Phone, &item.Notes, &item.CreatedAt, &subjects)
	if err != nil {
		return Tutor{}, err
	}
	item.Subjects = make([]string, 0)
	if subjects != "" {
		item.Subjects = strings.Split(subjects, ",")
	}
	return item, nil
}
