package config

import "os"

type Config struct {
	ServiceName string
	Host        string
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func Load(serviceName string, defaultPort string) Config {
	return Config{
		ServiceName: serviceName,
		Host:        getEnv("HOST", "127.0.0.1"),
		Port:        getEnv("PORT", defaultPort),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://repe_teacher:repe_teacher@localhost:5432/repe_teacher?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-repe-teacher-jwt-secret"),
	}
}

func getEnv(name string, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	return value
}
