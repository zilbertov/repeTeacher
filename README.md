Запуск проекта:

cd backend
docker compose up -d postgres

make migrate-up
make run-users

во втором терминале:

cd backend
make run-lessons

проверка состояния:

http://127.0.0.1:8081/health
http://127.0.0.1:8082/health

интерфейс:

frontend/index.html

Swagger:

http://127.0.0.1:8081/swagger/index.html
http://127.0.0.1:8082/swagger/index.html

Тесты:

cd backend
GOCACHE=/private/tmp/repeTeacher-go-build go test ./common/... ./users-service/... ./lessons-service/... -cover
