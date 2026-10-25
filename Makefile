.PHONY: run db migrate down

run:
	port=:3000 go run ./cmd/server/main.go

db:
	docker compose up -d db

migrate: db
	docker compose run --rm migrate

down:
	docker compose down

sqlc:
	