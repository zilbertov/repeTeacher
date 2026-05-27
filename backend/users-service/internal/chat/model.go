package chat

import "time"

type Chat struct {
	ID              int64     `json:"id"`
	TutorID         int64     `json:"tutor_id"`
	StudentID       int64     `json:"student_id"`
	ParticipantName string    `json:"participant_name"`
	LastMessage     string    `json:"last_message"`
	LastMessageTime time.Time `json:"last_message_time"`
}

type Message struct {
	ID         int64     `json:"id"`
	ChatID     int64     `json:"chat_id"`
	SenderType string    `json:"sender_type"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateChatRequest struct {
	StudentID int64 `json:"student_id"`
	TutorID   int64 `json:"tutor_id"`
}

type SendMessageRequest struct {
	Text       string `json:"text"`
	SenderType string `json:"sender_type"`
}
