package profile

import (
	"context"
	"database/sql"
	"strings"
)

type Repository interface {
	Get(ctx context.Context, tutorID int64) (Profile, error)
	Update(ctx context.Context, tutorID int64, req UpdateProfileRequest) (Profile, error)
	ChangePassword(ctx context.Context, tutorID int64, passwordHash string) error
	GetSettings(ctx context.Context, tutorID int64) (NotificationSettings, error)
	UpdateSettings(ctx context.Context, tutorID int64, settings NotificationSettings) (NotificationSettings, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Get(ctx context.Context, tutorID int64) (Profile, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT t.id, t.name, t.email, t.phone, COALESCE(string_agg(ts.subject, ',' ORDER BY ts.id), '')
		FROM tutors t
		LEFT JOIN tutor_subjects ts ON ts.tutor_id = t.id
		WHERE t.id = $1
		GROUP BY t.id
	`, tutorID)

	var item Profile
	var subjects string
	if err := row.Scan(&item.ID, &item.Name, &item.Email, &item.Phone, &subjects); err != nil {
		return Profile{}, err
	}
	if subjects != "" {
		item.Subjects = strings.Split(subjects, ",")
	}
	return item, nil
}

func (r *PostgresRepository) Update(ctx context.Context, tutorID int64, req UpdateProfileRequest) (Profile, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Profile{}, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE tutors
		SET name = $1, email = $2, phone = $3
		WHERE id = $4
	`, req.Name, req.Email, req.Phone, tutorID)
	if err != nil {
		return Profile{}, err
	}

	if req.Subjects != nil {
		if _, err = tx.ExecContext(ctx, `DELETE FROM tutor_subjects WHERE tutor_id = $1`, tutorID); err != nil {
			return Profile{}, err
		}
		for _, subject := range req.Subjects {
			if _, err = tx.ExecContext(ctx, `
				INSERT INTO tutor_subjects (tutor_id, subject)
				VALUES ($1, $2)
			`, tutorID, subject); err != nil {
				return Profile{}, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return Profile{}, err
	}
	return r.Get(ctx, tutorID)
}

func (r *PostgresRepository) ChangePassword(ctx context.Context, tutorID int64, passwordHash string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE tutors SET password_hash = $1 WHERE id = $2`, passwordHash, tutorID)
	return err
}

func (r *PostgresRepository) GetSettings(ctx context.Context, tutorID int64) (NotificationSettings, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT push_enabled, telegram_enabled, sound_enabled, lesson_reminders_enabled
		FROM notification_settings
		WHERE tutor_id = $1
	`, tutorID)

	var settings NotificationSettings
	err := row.Scan(
		&settings.PushEnabled,
		&settings.TelegramEnabled,
		&settings.SoundEnabled,
		&settings.LessonRemindersEnabled,
	)
	return settings, err
}

func (r *PostgresRepository) UpdateSettings(ctx context.Context, tutorID int64, settings NotificationSettings) (NotificationSettings, error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notification_settings
		SET push_enabled = $1,
		    telegram_enabled = $2,
		    sound_enabled = $3,
		    lesson_reminders_enabled = $4
		WHERE tutor_id = $5
	`, settings.PushEnabled, settings.TelegramEnabled, settings.SoundEnabled, settings.LessonRemindersEnabled, tutorID)
	if err != nil {
		return NotificationSettings{}, err
	}
	return r.GetSettings(ctx, tutorID)
}
