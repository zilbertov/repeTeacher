package app

import (
	"database/sql"
	"net/http"

	"github.com/swaggo/http-swagger"
	commonauth "github.com/zilbertov/repe-teacher-common/auth"
	"github.com/zilbertov/repe-teacher-common/response"
	"github.com/zilbertov/repe-teacher-lessons-service/internal/lesson"
)

func NewRouter(db *sql.DB, jwtSecret string) http.Handler {
	mux := http.NewServeMux()

	lessonRepo := lesson.NewPostgresRepository(db)
	lessonHandler := lesson.NewHandler(lesson.NewService(lessonRepo))

	mux.HandleFunc("/health", health)
	mux.Handle("/api/lessons", commonauth.Require(jwtSecret, http.HandlerFunc(lessonHandler.ListOrCreate)))
	mux.Handle("/api/lessons/", commonauth.Require(jwtSecret, http.HandlerFunc(lessonHandler.ByID)))
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}

func health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "lessons-service",
	})
}
