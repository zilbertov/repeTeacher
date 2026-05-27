package student

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

var ErrNotFound = errors.New("student not found")

type Repository interface {
	List(ctx context.Context, tutorID int64) ([]Student, error)
	Get(ctx context.Context, id int64) (Student, error)
	Create(ctx context.Context, tutorID int64, req CreateStudentRequest) (Student, error)
	Update(ctx context.Context, id int64, req UpdateStudentRequest) (Student, error)
	SetStatus(ctx context.Context, id int64, status string) (Student, error)
	UpdateNotes(ctx context.Context, id int64, notes string) (Student, error)
	Delete(ctx context.Context, id int64) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) List(ctx context.Context, tutorID int64) ([]Student, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.tutor_id, s.name, s.email, s.phone, s.exam_type, s.status, s.notes,
		       s.created_at, s.updated_at, COALESCE(string_agg(ss.subject, ','), '')
		FROM students s
		LEFT JOIN student_subjects ss ON ss.student_id = s.id
		WHERE s.tutor_id = $1
		GROUP BY s.id
		ORDER BY s.id
	`, tutorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := make([]Student, 0)
	for rows.Next() {
		item, err := scanStudent(rows)
		if err != nil {
			return nil, err
		}
		students = append(students, item)
	}

	return students, rows.Err()
}

func (r *PostgresRepository) Get(ctx context.Context, id int64) (Student, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT s.id, s.tutor_id, s.name, s.email, s.phone, s.exam_type, s.status, s.notes,
		       s.created_at, s.updated_at, COALESCE(string_agg(ss.subject, ','), '')
		FROM students s
		LEFT JOIN student_subjects ss ON ss.student_id = s.id
		WHERE s.id = $1
		GROUP BY s.id
	`, id)

	item, err := scanStudent(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Student{}, ErrNotFound
	}
	return item, err
}

func (r *PostgresRepository) Create(ctx context.Context, tutorID int64, req CreateStudentRequest) (Student, error) {
	status := req.Status
	if status == "" {
		status = "request"
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Student{}, err
	}
	defer tx.Rollback()

	var id int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO students (tutor_id, name, email, phone, exam_type, status, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, tutorID, req.Name, req.Email, req.Phone, req.ExamType, status, req.Notes).Scan(&id)
	if err != nil {
		return Student{}, err
	}

	for _, subject := range req.Subjects {
		if strings.TrimSpace(subject) == "" {
			continue
		}
		_, err = tx.ExecContext(ctx, `INSERT INTO student_subjects (student_id, subject) VALUES ($1, $2)`, id, subject)
		if err != nil {
			return Student{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return Student{}, err
	}

	return r.Get(ctx, id)
}

func (r *PostgresRepository) Update(ctx context.Context, id int64, req UpdateStudentRequest) (Student, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Student{}, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		UPDATE students
		SET name = $1, email = $2, phone = $3, exam_type = $4, status = $5, notes = $6, updated_at = now()
		WHERE id = $7
	`, req.Name, req.Email, req.Phone, req.ExamType, req.Status, req.Notes, id)
	if err != nil {
		return Student{}, err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return Student{}, ErrNotFound
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM student_subjects WHERE student_id = $1`, id)
	if err != nil {
		return Student{}, err
	}

	for _, subject := range req.Subjects {
		if strings.TrimSpace(subject) == "" {
			continue
		}
		_, err = tx.ExecContext(ctx, `INSERT INTO student_subjects (student_id, subject) VALUES ($1, $2)`, id, subject)
		if err != nil {
			return Student{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return Student{}, err
	}

	return r.Get(ctx, id)
}

func (r *PostgresRepository) SetStatus(ctx context.Context, id int64, status string) (Student, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE students
		SET status = $1, updated_at = now()
		WHERE id = $2
	`, status, id)
	if err != nil {
		return Student{}, err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return Student{}, ErrNotFound
	}

	return r.Get(ctx, id)
}

func (r *PostgresRepository) UpdateNotes(ctx context.Context, id int64, notes string) (Student, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE students
		SET notes = $1, updated_at = now()
		WHERE id = $2
	`, notes, id)
	if err != nil {
		return Student{}, err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return Student{}, ErrNotFound
	}

	return r.Get(ctx, id)
}

func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM students WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return ErrNotFound
	}
	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanStudent(row scanner) (Student, error) {
	var item Student
	var subjects string
	err := row.Scan(
		&item.ID,
		&item.TutorID,
		&item.Name,
		&item.Email,
		&item.Phone,
		&item.ExamType,
		&item.Status,
		&item.Notes,
		&item.CreatedAt,
		&item.UpdatedAt,
		&subjects,
	)
	if err != nil {
		return Student{}, err
	}
	item.Subjects = make([]string, 0)
	if subjects != "" {
		item.Subjects = strings.Split(subjects, ",")
	}
	return item, nil
}
