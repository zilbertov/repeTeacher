package profile

import (
	"context"
	"errors"
	"testing"
)

type fakeProfileRepo struct {
	profile  Profile
	settings NotificationSettings
	password string
}

func newFakeProfileRepo() *fakeProfileRepo {
	return &fakeProfileRepo{
		profile: Profile{ID: 1, Name: "Вадим", Email: "v4bem@ya.ru", Phone: "89198318673", Subjects: []string{"Математика"}},
		settings: NotificationSettings{
			PushEnabled:            true,
			TelegramEnabled:        true,
			SoundEnabled:           true,
			LessonRemindersEnabled: false,
		},
	}
}

func (r *fakeProfileRepo) Get(ctx context.Context, tutorID int64) (Profile, error) {
	return r.profile, nil
}

func (r *fakeProfileRepo) Update(ctx context.Context, tutorID int64, req UpdateProfileRequest) (Profile, error) {
	r.profile.Name = req.Name
	r.profile.Email = req.Email
	r.profile.Phone = req.Phone
	if req.Subjects != nil {
		r.profile.Subjects = req.Subjects
	}
	return r.profile, nil
}

func (r *fakeProfileRepo) ChangePassword(ctx context.Context, tutorID int64, passwordHash string) error {
	r.password = passwordHash
	return nil
}

func (r *fakeProfileRepo) GetSettings(ctx context.Context, tutorID int64) (NotificationSettings, error) {
	return r.settings, nil
}

func (r *fakeProfileRepo) UpdateSettings(ctx context.Context, tutorID int64, settings NotificationSettings) (NotificationSettings, error) {
	r.settings = settings
	return r.settings, nil
}

func TestUpdateProfile(t *testing.T) {
	service := NewService(newFakeProfileRepo())

	item, err := service.Update(context.Background(), 1, UpdateProfileRequest{Name: "Новое имя", Email: "new@email.ru", Phone: "123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Name != "Новое имя" {
		t.Fatalf("profile was not updated")
	}
}

func TestUpdateProfileNeedsNameAndEmail(t *testing.T) {
	service := NewService(newFakeProfileRepo())

	_, err := service.Update(context.Background(), 1, UpdateProfileRequest{Name: "", Email: ""})
	if !errors.Is(err, ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestUpdateProfileCleansSubjects(t *testing.T) {
	service := NewService(newFakeProfileRepo())

	item, err := service.Update(context.Background(), 1, UpdateProfileRequest{
		Name:     "Вадим",
		Email:    "v4bem@ya.ru",
		Phone:    "89198318673",
		Subjects: []string{" Математика ", "", "Математика", "Русский язык"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(item.Subjects) != 2 {
		t.Fatalf("expected 2 subjects, got %d", len(item.Subjects))
	}
	if item.Subjects[0] != "Математика" || item.Subjects[1] != "Русский язык" {
		t.Fatalf("subjects were not cleaned: %#v", item.Subjects)
	}
}

func TestChangePasswordNeedsLongEnoughPassword(t *testing.T) {
	service := NewService(newFakeProfileRepo())

	err := service.ChangePassword(context.Background(), 1, ChangePasswordRequest{NewPassword: "123"})
	if !errors.Is(err, ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}
