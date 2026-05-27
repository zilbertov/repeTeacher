package main

import (
	"context"
	"net"
	"net/http"

	_ "github.com/lib/pq"
	_ "github.com/zilbertov/repe-teacher-users-service/docs"

	"github.com/zilbertov/repe-teacher-common/config"
	"github.com/zilbertov/repe-teacher-common/db"
	"github.com/zilbertov/repe-teacher-common/logger"
	"github.com/zilbertov/repe-teacher-common/middleware"
	"github.com/zilbertov/repe-teacher-users-service/internal/app"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load("users-service", "8081")

	log, err := logger.New(cfg.ServiceName)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	ctx := context.Background()
	conn, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("database connection failed", zap.Error(err))
	}
	defer conn.Close()

	router := app.NewRouter(conn, cfg.JWTSecret)
	handler := middleware.CORS(middleware.RequestLogger(log, router))

	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	log.Info("users-service started", zap.String("addr", addr))
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal("server stopped", zap.Error(err))
	}
}
