start-docker:
	docker compose up -d
migrate-create:
	migrate create -ext sql -dir db/migration -seq init_schema
migrate-up:
	migrate -path db/migration -database "postgres://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrate-down:
	migrate -path db/migration -database "postgres://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down

.PHONY: start-docker migrate-create migrate-up migrate-down
