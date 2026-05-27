package lesson

import (
	"context"
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("lesson not found")

type Repository interface {
	List(ctx context.Context, tutorID int64) ([]Lesson, error)
	ListByStudent(ctx context.Context, studentID int64) ([]Lesson, error)
	Get(ctx context.Context, id int64) (Lesson, error)
	Create(ctx context.Context, tutorID int64, req CreateLessonRequest) (Lesson, error)
	Update(ctx context.Context, id int64, req UpdateLessonRequest) (Lesson, error)
	Reschedule(ctx context.Context, id int64, req RescheduleRequest) (Lesson, error)
	Cancel(ctx context.Context, id int64, req CancelRequest) (Lesson, error)
	AddFile(ctx context.Context, id int64, req AddFileRequest) (LessonFile, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) List(ctx context.Context, tutorID int64) ([]Lesson, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT l.id, l.tutor_id, l.student_id, t.name, s.name, l.subject, l.exam_type,
		       l.lesson_date::text, l.start_time::text, l.duration_minutes, l.format,
		       l.has_homework, l.price, l.status, l.created_at, l.updated_at
		FROM lessons l
		JOIN tutors t ON t.id = l.tutor_id
		JOIN students s ON s.id = l.student_id
		WHERE l.tutor_id = $1
		ORDER BY l.lesson_date, l.start_time
	`, tutorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Lesson, 0)
	for rows.Next() {
		item, err := scanLesson(rows)
		if err != nil {
			return nil, err
		}
		files, err := r.listFiles(ctx, item.ID)
		if err != nil {
			return nil, err
		}
		item.Files = files
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) ListByStudent(ctx context.Context, studentID int64) ([]Lesson, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT l.id, l.tutor_id, l.student_id, t.name, s.name, l.subject, l.exam_type,
		       l.lesson_date::text, l.start_time::text, l.duration_minutes, l.format,
		       l.has_homework, l.price, l.status, l.created_at, l.updated_at
		FROM lessons l
		JOIN tutors t ON t.id = l.tutor_id
		JOIN students s ON s.id = l.student_id
		WHERE l.student_id = $1
		ORDER BY l.lesson_date, l.start_time
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Lesson, 0)
	for rows.Next() {
		item, err := scanLesson(rows)
		if err != nil {
			return nil, err
		}
		files, err := r.listFiles(ctx, item.ID)
		if err != nil {
			return nil, err
		}
		item.Files = files
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) Get(ctx context.Context, id int64) (Lesson, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT l.id, l.tutor_id, l.student_id, t.name, s.name, l.subject, l.exam_type,
		       l.lesson_date::text, l.start_time::text, l.duration_minutes, l.format,
		       l.has_homework, l.price, l.status, l.created_at, l.updated_at
		FROM lessons l
		JOIN tutors t ON t.id = l.tutor_id
		JOIN students s ON s.id = l.student_id
		WHERE l.id = $1
	`, id)

	item, err := scanLesson(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Lesson{}, ErrNotFound
	}
	if err != nil {
		return Lesson{}, err
	}

	files, err := r.listFiles(ctx, item.ID)
	if err != nil {
		return Lesson{}, err
	}
	item.Files = files
	return item, nil
}

func (r *PostgresRepository) Create(ctx context.Context, tutorID int64, req CreateLessonRequest) (Lesson, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Lesson{}, err
	}
	defer tx.Rollback()

	var id int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO lessons (tutor_id, student_id, subject, exam_type, lesson_date, start_time,
		                     duration_minutes, format, has_homework, price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`, tutorID, req.StudentID, req.Subject, req.ExamType, req.LessonDate, req.StartTime,
		req.DurationMinutes, req.Format, req.HasHomework, req.Price).Scan(&id)
	if err != nil {
		return Lesson{}, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO notifications (tutor_id, student_id, lesson_id, type, title, description, recipient_type)
		VALUES ($1, $2, $3, 'lesson', 'Новое занятие', 'Репетитор добавил новое занятие.', 'student')
	`, tutorID, req.StudentID, id)
	if err != nil {
		return Lesson{}, err
	}

	if err := tx.Commit(); err != nil {
		return Lesson{}, err
	}
	return r.Get(ctx, id)
}

func (r *PostgresRepository) Update(ctx context.Context, id int64, req UpdateLessonRequest) (Lesson, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE lessons
		SET student_id = $1, subject = $2, exam_type = $3, lesson_date = $4, start_time = $5,
		    duration_minutes = $6, format = $7, has_homework = $8, price = $9, updated_at = now()
		WHERE id = $10
	`, req.StudentID, req.Subject, req.ExamType, req.LessonDate, req.StartTime,
		req.DurationMinutes, req.Format, req.HasHomework, req.Price, id)
	if err != nil {
		return Lesson{}, err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return Lesson{}, ErrNotFound
	}
	return r.Get(ctx, id)
}

func (r *PostgresRepository) Reschedule(ctx context.Context, id int64, req RescheduleRequest) (Lesson, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Lesson{}, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		UPDATE lessons
		SET lesson_date = $1, start_time = $2, updated_at = now()
		WHERE id = $3
	`, req.LessonDate, req.StartTime, id)
	if err != nil {
		return Lesson{}, err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return Lesson{}, ErrNotFound
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO notifications (tutor_id, student_id, lesson_id, type, title, description, recipient_type)
		SELECT tutor_id, student_id, id, 'reschedule', 'Перенос занятия', $2, $3
		FROM lessons
		WHERE id = $1
	`, id, notificationText("reschedule", req.SenderType), recipientType(req.SenderType))
	if err != nil {
		return Lesson{}, err
	}

	if err := tx.Commit(); err != nil {
		return Lesson{}, err
	}
	return r.Get(ctx, id)
}

func (r *PostgresRepository) Cancel(ctx context.Context, id int64, req CancelRequest) (Lesson, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Lesson{}, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		UPDATE lessons
		SET status = 'cancelled', updated_at = now()
		WHERE id = $1
	`, id)
	if err != nil {
		return Lesson{}, err
	}
	if changed, _ := result.RowsAffected(); changed == 0 {
		return Lesson{}, ErrNotFound
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO notifications (tutor_id, student_id, lesson_id, type, title, description, recipient_type)
		SELECT tutor_id, student_id, id, 'cancel', 'Отмена занятия', $2, $3
		FROM lessons
		WHERE id = $1
	`, id, notificationText("cancel", req.SenderType), recipientType(req.SenderType))
	if err != nil {
		return Lesson{}, err
	}

	if err := tx.Commit(); err != nil {
		return Lesson{}, err
	}
	return r.Get(ctx, id)
}

func (r *PostgresRepository) AddFile(ctx context.Context, id int64, req AddFileRequest) (LessonFile, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO lesson_files (lesson_id, file_type, file_name, file_path)
		VALUES ($1, $2, $3, $4)
		RETURNING id, lesson_id, file_type, file_name, file_path
	`, id, req.FileType, req.FileName, req.FilePath)

	var file LessonFile
	err := row.Scan(&file.ID, &file.LessonID, &file.FileType, &file.FileName, &file.FilePath)
	return file, err
}

func (r *PostgresRepository) listFiles(ctx context.Context, lessonID int64) ([]LessonFile, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, lesson_id, file_type, file_name, file_path
		FROM lesson_files
		WHERE lesson_id = $1
		ORDER BY id
	`, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]LessonFile, 0)
	for rows.Next() {
		var file LessonFile
		err := rows.Scan(&file.ID, &file.LessonID, &file.FileType, &file.FileName, &file.FilePath)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, rows.Err()
}

type scanner interface {
	Scan(dest ...any) error
}

func scanLesson(row scanner) (Lesson, error) {
	var item Lesson
	err := row.Scan(
		&item.ID,
		&item.TutorID,
		&item.StudentID,
		&item.TutorName,
		&item.StudentName,
		&item.Subject,
		&item.ExamType,
		&item.LessonDate,
		&item.StartTime,
		&item.DurationMinutes,
		&item.Format,
		&item.HasHomework,
		&item.Price,
		&item.Status,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	return item, err
}

func recipientType(senderType string) string {
	if senderType == "student" {
		return "tutor"
	}
	return "student"
}

func notificationText(eventType string, senderType string) string {
	if eventType == "cancel" {
		if senderType == "student" {
			return "Ученик отменил занятие."
		}
		return "Занятие отменено репетитором."
	}
	if senderType == "student" {
		return "Ученик запросил перенос занятия."
	}
	return "Занятие перенесено репетитором."
}
