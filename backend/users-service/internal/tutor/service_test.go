package tutor

import (
	"context"
	"errors"
	"testing"
)

type fakeTutorRepo struct {
	items []Tutor
}

func (r *fakeTutorRepo) List(ctx context.Context) ([]Tutor, error) {
	return r.items, nil
}

func (r *fakeTutorRepo) Get(ctx context.Context, id int64) (Tutor, error) {
	for _, item := range r.items {
		if item.ID == id {
			return item, nil
		}
	}
	return Tutor{}, ErrNotFound
}

func (r *fakeTutorRepo) Create(ctx context.Context, req CreateTutorRequest) (Tutor, error) {
	item := Tutor{ID: int64(len(r.items) + 1), Name: req.Name, Email: req.Email, Phone: req.Phone, Subjects: req.Subjects, Notes: req.Notes}
	r.items = append(r.items, item)
	return item, nil
}

func (r *fakeTutorRepo) UpdateNotes(ctx context.Context, id int64, notes string) (Tutor, error) {
	for index, item := range r.items {
		if item.ID == id {
			item.Notes = notes
			r.items[index] = item
			return item, nil
		}
	}
	return Tutor{}, ErrNotFound
}

func TestCreateTutor(t *testing.T) {
	service := NewService(&fakeTutorRepo{})

	item, err := service.Create(context.Background(), CreateTutorRequest{
		Name:     "Анна Смирнова",
		Email:    "anna@example.com",
		Phone:    "89191234567",
		Subjects: []string{"Математика", "Математика"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Name != "Анна Смирнова" || len(item.Subjects) != 1 {
		t.Fatalf("unexpected tutor: %+v", item)
	}
}

func TestCreateTutorValidation(t *testing.T) {
	service := NewService(&fakeTutorRepo{})

	_, err := service.Create(context.Background(), CreateTutorRequest{Name: "", Email: ""})
	if !errors.Is(err, ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestUpdateTutorNotes(t *testing.T) {
	service := NewService(&fakeTutorRepo{items: []Tutor{{ID: 1, Name: "Анна", Email: "anna@example.com", Subjects: []string{"Математика"}}}})

	item, err := service.UpdateNotes(context.Background(), 1, UpdateTutorNotesRequest{Notes: " Хорошо объясняет "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Notes != "Хорошо объясняет" {
		t.Fatalf("notes were not trimmed and saved: %q", item.Notes)
	}
}
