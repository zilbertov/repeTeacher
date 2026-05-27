package session

type LoginRequest struct {
	Role     string `json:"role"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	Role      string `json:"role"`
	TutorID   int64  `json:"tutor_id,omitempty"`
	StudentID int64  `json:"student_id,omitempty"`
	ExpiresAt int64  `json:"expires_at"`
}
