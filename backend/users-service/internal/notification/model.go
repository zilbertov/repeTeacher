package notification

import "time"

type Notification struct {
	ID            int64     `json:"id"`
	TutorID       int64     `json:"tutor_id"`
	StudentID     *int64    `json:"student_id"`
	LessonID      *int64    `json:"lesson_id"`
	Type          string    `json:"type"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	RecipientType string    `json:"recipient_type"`
	IsRead        bool      `json:"is_read"`
	CreatedAt     time.Time `json:"created_at"`
}
