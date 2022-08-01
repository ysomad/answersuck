include .env
export

MIGRATE := migrate -path migrations -database "$(PG_URL)?sslmode=disable"
DB_SOURCE=postgresql://user:pass@localhost:5432/postgres?sslmode=disable
.PHONY: compose-up
compose-up:
	docker compose up --build -d postgres && docker compose logs -f

.PHONY: compose-down
compose-down:
	docker-compose down --remove-orphans

.PHONY: run
run:
	go mod tidy && go mod download && \
	go run ./cmd/app

.PHONY: build
build:
	go build ./cmd/app

.PHONY: test
test:
	go test -v -cover -race -count 1 ./internal/... ./pkg/...

.PHONY: compose-up-integration-test
compose-up-integration-test:
	docker-compose up --build --abort-on-container-exit --exit-code-from integration

.PHONY: integration-test
integration-test:
	LOG_LEVEL=debug go test -v -race -count 1 ./test/...

.PHONY: migrate-new
migrate-new:
	@read -p "Enter the name of the new migration: " name; \
	$(MIGRATE) create -ext sql -dir migrations $${name// /_}

.PHONY: migrate-up
migrate-up:
	@echo "Running all new database migrations..."
	@$(MIGRATE) up

.PHONY: migrate-drop
migrate-drop:
	@echo "Dropping everything in database..."
	@$(MIGRATE) drop

.PHONY: migrate-down
migrate-down:
	@echo "Running all down database migrations..."
	@$(MIGRATE) down
migrateup:
	migrate -path migrations -database "$(DB_SOURCE)" -verbose up

