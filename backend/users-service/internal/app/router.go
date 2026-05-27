package app

import (
	"database/sql"
	"net/http"

	"github.com/swaggo/http-swagger"
	commonauth "github.com/zilbertov/repe-teacher-common/auth"
	"github.com/zilbertov/repe-teacher-common/response"
	"github.com/zilbertov/repe-teacher-users-service/internal/chat"
	"github.com/zilbertov/repe-teacher-users-service/internal/notification"
	"github.com/zilbertov/repe-teacher-users-service/internal/profile"
	"github.com/zilbertov/repe-teacher-users-service/internal/session"
	"github.com/zilbertov/repe-teacher-users-service/internal/student"
	"github.com/zilbertov/repe-teacher-users-service/internal/tutor"
)

func NewRouter(db *sql.DB, jwtSecret string) http.Handler {
	mux := http.NewServeMux()

	sessionRepo := session.NewPostgresRepository(db)
	sessionHandler := session.NewHandler(session.NewService(sessionRepo, jwtSecret))

	studentRepo := student.NewPostgresRepository(db)
	studentHandler := student.NewHandler(student.NewService(studentRepo))

	profileRepo := profile.NewPostgresRepository(db)
	profileHandler := profile.NewHandler(profile.NewService(profileRepo))

	chatRepo := chat.NewPostgresRepository(db)
	chatHandler := chat.NewHandler(chat.NewService(chatRepo))

	notificationRepo := notification.NewPostgresRepository(db)
	notificationHandler := notification.NewHandler(notification.NewService(notificationRepo))

	tutorRepo := tutor.NewPostgresRepository(db)
	tutorHandler := tutor.NewHandler(tutor.NewService(tutorRepo))

	mux.HandleFunc("/health", health)
	mux.HandleFunc("/api/auth/login", sessionHandler.Login)
	mux.Handle("/api/students", commonauth.Require(jwtSecret, http.HandlerFunc(studentHandler.ListOrCreate)))
	mux.Handle("/api/students/", commonauth.Require(jwtSecret, http.HandlerFunc(studentHandler.ByID)))
	mux.Handle("/api/tutors", commonauth.Require(jwtSecret, http.HandlerFunc(tutorHandler.ListOrCreate)))
	mux.Handle("/api/tutors/", commonauth.Require(jwtSecret, http.HandlerFunc(tutorHandler.ByID)))
	mux.Handle("/api/profile", commonauth.Require(jwtSecret, http.HandlerFunc(profileHandler.Profile)))
	mux.Handle("/api/profile/password", commonauth.Require(jwtSecret, http.HandlerFunc(profileHandler.Password)))
	mux.Handle("/api/settings/notifications", commonauth.Require(jwtSecret, http.HandlerFunc(profileHandler.Settings)))
	mux.Handle("/api/chats", commonauth.Require(jwtSecret, http.HandlerFunc(chatHandler.ListOrCreate)))
	mux.Handle("/api/chats/", commonauth.Require(jwtSecret, http.HandlerFunc(chatHandler.Messages)))
	mux.Handle("/api/notifications", commonauth.Require(jwtSecret, http.HandlerFunc(notificationHandler.List)))
	mux.Handle("/api/notifications/", commonauth.Require(jwtSecret, http.HandlerFunc(notificationHandler.ByID)))
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}

func health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "users-service",
	})
}
