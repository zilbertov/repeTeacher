package student

import "time"

type Student struct {
	ID        int64     `json:"id"`
	TutorID   int64     `json:"tutor_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Subjects  []string  `json:"subjects"`
	ExamType  string    `json:"exam_type"`
	Status    string    `json:"status"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateStudentRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Subjects []string `json:"subjects"`
	ExamType string   `json:"exam_type"`
	Status   string   `json:"status"`
	Notes    string   `json:"notes"`
}

type UpdateStudentRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Subjects []string `json:"subjects"`
	ExamType string   `json:"exam_type"`
	Status   string   `json:"status"`
	Notes    string   `json:"notes"`
}

type NotesRequest struct {
	Notes string `json:"notes"`
}
