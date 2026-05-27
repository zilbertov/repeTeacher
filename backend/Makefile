DATABASE_URL ?= postgres://repe_teacher:repe_teacher@127.0.0.1:5432/repe_teacher?sslmode=disable

.PHONY: postgres migrate-up migrate-down run-users run-lessons test

postgres:
	docker compose up -d postgres

migrate-up:
	DATABASE_URL="$(DATABASE_URL)" go run ./tools/migrate up

migrate-down:
	DATABASE_URL="$(DATABASE_URL)" go run ./tools/migrate down

migrate-force:
	DATABASE_URL="$(DATABASE_URL)" go run ./tools/migrate force $(VERSION)

reset-db:
	docker compose down -v

run-users:
	cd users-service && PORT=8081 DATABASE_URL="$(DATABASE_URL)" go run ./cmd/users-service

run-lessons:
	cd lessons-service && PORT=8082 DATABASE_URL="$(DATABASE_URL)" go run ./cmd/lessons-service

test:
	go test ./common/... ./users-service/... ./lessons-service/... -cover
