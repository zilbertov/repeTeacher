package profile

type Profile struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Subjects []string `json:"subjects"`
}

type UpdateProfileRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Subjects []string `json:"subjects"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type NotificationSettings struct {
	PushEnabled            bool `json:"push_enabled"`
	TelegramEnabled        bool `json:"telegram_enabled"`
	SoundEnabled           bool `json:"sound_enabled"`
	LessonRemindersEnabled bool `json:"lesson_reminders_enabled"`
}
