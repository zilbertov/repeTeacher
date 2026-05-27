package tutor

import "time"

type Tutor struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Subjects  []string  `json:"subjects"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTutorRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Subjects []string `json:"subjects"`
	Notes    string   `json:"notes"`
}

type UpdateTutorNotesRequest struct {
	Notes string `json:"notes"`
}
