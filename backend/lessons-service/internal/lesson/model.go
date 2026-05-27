package lesson

import "time"

type Lesson struct {
	ID              int64        `json:"id"`
	TutorID         int64        `json:"tutor_id"`
	StudentID       int64        `json:"student_id"`
	TutorName       string       `json:"tutor_name"`
	StudentName     string       `json:"student_name"`
	Subject         string       `json:"subject"`
	ExamType        string       `json:"exam_type"`
	LessonDate      string       `json:"lesson_date"`
	StartTime       string       `json:"start_time"`
	DurationMinutes int          `json:"duration_minutes"`
	Format          string       `json:"format"`
	HasHomework     bool         `json:"has_homework"`
	Price           int          `json:"price"`
	Status          string       `json:"status"`
	Files           []LessonFile `json:"files"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type LessonFile struct {
	ID       int64  `json:"id"`
	LessonID int64  `json:"lesson_id"`
	FileType string `json:"file_type"`
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
}

type CreateLessonRequest struct {
	TutorID         int64  `json:"tutor_id"`
	StudentID       int64  `json:"student_id"`
	Subject         string `json:"subject"`
	ExamType        string `json:"exam_type"`
	LessonDate      string `json:"lesson_date"`
	StartTime       string `json:"start_time"`
	DurationMinutes int    `json:"duration_minutes"`
	Format          string `json:"format"`
	HasHomework     bool   `json:"has_homework"`
	Price           int    `json:"price"`
}

type UpdateLessonRequest = CreateLessonRequest

type RescheduleRequest struct {
	LessonDate string `json:"lesson_date"`
	StartTime  string `json:"start_time"`
	SenderType string `json:"sender_type"`
}

type CancelRequest struct {
	SenderType string `json:"sender_type"`
}

type AddFileRequest struct {
	FileType string `json:"file_type"`
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
}
