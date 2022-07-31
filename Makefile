include .env
export

MIGRATE := migrate -path migrations -database "$(PG_URL)?sslmode=disable"

.PHONY: compose-up
compose-up:
	docker compose up --build -d postgres redis && docker compose logs -f

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
	@$(MIGRATE) down

.PHONY: test
test:
	INTEGRATION_TESTDB=true INTEGRATION_LOGLEVEL=debug go test -v -race -count=1 ./...
