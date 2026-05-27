package student

import (
	"context"
	"errors"
	"testing"
)

type fakeStudentRepo struct {
	items map[int64]Student
}

func newFakeStudentRepo() *fakeStudentRepo {
	return &fakeStudentRepo{items: map[int64]Student{
		1: {ID: 1, Name: "Илья", Email: "ilya@email.ru", Subjects: []string{"Математика"}, Status: "request"},
		2: {ID: 2, Name: "Полина", Email: "polina@email.ru", Subjects: []string{"Русский язык"}, Status: "active"},
	}}
}

func (r *fakeStudentRepo) List(ctx context.Context, tutorID int64) ([]Student, error) {
	return []Student{r.items[1], r.items[2]}, nil
}

func (r *fakeStudentRepo) Get(ctx context.Context, id int64) (Student, error) {
	item, ok := r.items[id]
	if !ok {
		return Student{}, ErrNotFound
	}
	return item, nil
}

func (r *fakeStudentRepo) Create(ctx context.Context, tutorID int64, req CreateStudentRequest) (Student, error) {
	item := Student{ID: 3, TutorID: tutorID, Name: req.Name, Email: req.Email, Subjects: req.Subjects, Status: req.Status}
	r.items[item.ID] = item
	return item, nil
}

func (r *fakeStudentRepo) Update(ctx context.Context, id int64, req UpdateStudentRequest) (Student, error) {
	item := r.items[id]
	item.Name = req.Name
	item.Email = req.Email
	item.Subjects = req.Subjects
	item.Status = req.Status
	r.items[id] = item
	return item, nil
}

func (r *fakeStudentRepo) SetStatus(ctx context.Context, id int64, status string) (Student, error) {
	item, ok := r.items[id]
	if !ok {
		return Student{}, ErrNotFound
	}
	item.Status = status
	r.items[id] = item
	return item, nil
}

func (r *fakeStudentRepo) UpdateNotes(ctx context.Context, id int64, notes string) (Student, error) {
	item := r.items[id]
	item.Notes = notes
	r.items[id] = item
	return item, nil
}

func (r *fakeStudentRepo) Delete(ctx context.Context, id int64) error {
	delete(r.items, id)
	return nil
}

func TestAcceptRequestStudent(t *testing.T) {
	service := NewService(newFakeStudentRepo())

	item, err := service.Accept(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Status != "active" {
		t.Fatalf("expected active status, got %s", item.Status)
	}
}

func TestAcceptActiveStudentReturnsError(t *testing.T) {
	service := NewService(newFakeStudentRepo())

	_, err := service.Accept(context.Background(), 2)
	if !errors.Is(err, ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestCreateStudentNeedsNameEmailAndSubject(t *testing.T) {
	service := NewService(newFakeStudentRepo())

	_, err := service.Create(context.Background(), 1, CreateStudentRequest{Name: "", Email: "", Subjects: nil})
	if !errors.Is(err, ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestUpdateNotes(t *testing.T) {
	service := NewService(newFakeStudentRepo())

	item, err := service.UpdateNotes(context.Background(), 2, "Нужна практика")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Notes != "Нужна практика" {
		t.Fatalf("expected notes to be changed")
	}
}

func TestArchiveStudent(t *testing.T) {
	service := NewService(newFakeStudentRepo())

	item, err := service.Archive(context.Background(), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Status != "archived" {
		t.Fatalf("expected archived status")
	}
}
