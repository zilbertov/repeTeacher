package session

import (
	"context"
	"errors"
	"strings"
	"time"

	commonauth "github.com/zilbertov/repe-teacher-common/auth"
)

var (
	ErrBadRequest         = errors.New("bad request")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

const demoPassword = "demo"

type Service struct {
	repo      Repository
	jwtSecret string
	tokenTTL  time.Duration
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: jwtSecret,
		tokenTTL:  8 * time.Hour,
	}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	role := strings.TrimSpace(req.Role)
	email := strings.TrimSpace(req.Email)
	if role == "" || email == "" {
		return LoginResponse{}, ErrBadRequest
	}
	if req.Password != demoPassword {
		return LoginResponse{}, ErrInvalidCredentials
	}

	claims := commonauth.Claims{Role: role}
	switch role {
	case commonauth.RoleTutor:
		tutorID, err := s.repo.FindTutorByEmail(ctx, email)
		if err != nil {
			return LoginResponse{}, err
		}
		claims.TutorID = tutorID
	case commonauth.RoleStudent:
		studentID, tutorID, err := s.repo.FindStudentByEmail(ctx, email)
		if err != nil {
			return LoginResponse{}, err
		}
		claims.StudentID = studentID
		claims.TutorID = tutorID
	default:
		return LoginResponse{}, ErrBadRequest
	}

	token, signedClaims, err := commonauth.NewToken(s.jwtSecret, claims, s.tokenTTL)
	if err != nil {
		return LoginResponse{}, err
	}
	return LoginResponse{
		Token:     token,
		Role:      signedClaims.Role,
		TutorID:   signedClaims.TutorID,
		StudentID: signedClaims.StudentID,
		ExpiresAt: signedClaims.ExpiresAt,
	}, nil
}
