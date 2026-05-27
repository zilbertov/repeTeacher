package auth

import (
	"testing"
	"time"
)

func TestNewAndParseTutorToken(t *testing.T) {
	token, _, err := NewToken("secret", Claims{Role: RoleTutor, TutorID: 1}, time.Hour)
	if err != nil {
		t.Fatalf("new token: %v", err)
	}

	claims, err := ParseToken("secret", token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if claims.Role != RoleTutor || claims.TutorID != 1 {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestParseRejectsWrongSecret(t *testing.T) {
	token, _, err := NewToken("secret", Claims{Role: RoleStudent, TutorID: 1, StudentID: 4}, time.Hour)
	if err != nil {
		t.Fatalf("new token: %v", err)
	}

	if _, err := ParseToken("other-secret", token); err != ErrUnauthorized {
		t.Fatalf("expected unauthorized, got %v", err)
	}
}
