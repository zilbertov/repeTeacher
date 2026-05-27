package profile

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepositoryGetProfileAndSettings(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT t.id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "phone", "subjects"}).
			AddRow(1, "Вадим", "v4bem@ya.ru", "89198318673", "Математика,Русский язык"))
	mock.ExpectQuery("SELECT push_enabled").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(settingsColumns()).AddRow(true, true, false, true))

	repo := NewPostgresRepository(db)
	item, err := repo.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("get profile: %v", err)
	}
	if len(item.Subjects) != 2 {
		t.Fatalf("unexpected subjects: %+v", item.Subjects)
	}

	settings, err := repo.GetSettings(context.Background(), 1)
	if err != nil {
		t.Fatalf("get settings: %v", err)
	}
	if !settings.PushEnabled || settings.SoundEnabled {
		t.Fatalf("unexpected settings: %+v", settings)
	}
}

func TestRepositoryUpdateProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tutors").
		WithArgs("Вадим", "new@email.ru", "123", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM tutor_subjects").WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO tutor_subjects").WithArgs(int64(1), "Математика").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO tutor_subjects").WithArgs(int64(1), "Русский язык").
		WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT t.id").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "phone", "subjects"}).
			AddRow(1, "Вадим", "new@email.ru", "123", "Математика,Русский язык"))

	repo := NewPostgresRepository(db)
	item, err := repo.Update(context.Background(), 1, UpdateProfileRequest{
		Name:     "Вадим",
		Email:    "new@email.ru",
		Phone:    "123",
		Subjects: []string{"Математика", "Русский язык"},
	})
	if err != nil {
		t.Fatalf("update profile: %v", err)
	}
	if item.Email != "new@email.ru" {
		t.Fatalf("email was not updated")
	}
}

func TestRepositoryPasswordAndSettingsUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE tutors SET password_hash").
		WithArgs("changed:1234", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE notification_settings").
		WithArgs(false, true, false, true, int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT push_enabled").WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(settingsColumns()).AddRow(false, true, false, true))

	repo := NewPostgresRepository(db)
	if err := repo.ChangePassword(context.Background(), 1, "changed:1234"); err != nil {
		t.Fatalf("change password: %v", err)
	}

	settings, err := repo.UpdateSettings(context.Background(), 1, NotificationSettings{
		PushEnabled:            false,
		TelegramEnabled:        true,
		SoundEnabled:           false,
		LessonRemindersEnabled: true,
	})
	if err != nil {
		t.Fatalf("update settings: %v", err)
	}
	if settings.PushEnabled || !settings.LessonRemindersEnabled {
		t.Fatalf("unexpected settings: %+v", settings)
	}
}

func settingsColumns() []string {
	return []string{"push_enabled", "telegram_enabled", "sound_enabled", "lesson_reminders_enabled"}
}
