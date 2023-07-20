include .env
export

MIGRATE := migrate -path migrations -database "$(PG_URL)?sslmode=disable"

.PHONY: compose-dev
compose-dev:
	docker-compose up --build -d postgres && docker-compose logs -f

.PHONY: compose-up
compose-up:
	docker-compose up --build -d postgres jaeger prometheus grafana && docker-compose logs -f

.PHONY: compose-all
compose-all:
	docker-compose up --build -d && docker-compose logs -f

.PHONY: compose-min
compose-min:
	docker-compose up --build -d postgres app && docker-compose logs -f

.PHONY: compose-down
compose-down:
	docker-compose down --remove-orphans

.PHONY: run
run:
	go mod tidy && go mod download && \
	go run ./cmd/app

.PHONY: run-migrate
run-migrate:
	go mod tidy && go mod download && \
	go run -tags migrate ./cmd/app

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -v -cover -race -count 1 ./internal/...

.PHONY: migrate-new
migrate-new:
	@read -p "Enter the name of the new migration: " name; \
	$(MIGRATE) create -ext sql -dir migrations $${name// /_}

.PHONY: migrate-up
migrate-up:
	@echo "Running all new database migrations..."
	@$(MIGRATE) up

.PHONY: migrate-down
migrate-down:
	@echo "Running all down database migrations..."
	@$(MIGRATE) down

.PHONY: migrate-drop
migrate-drop:
	@echo "Dropping everything in database..."
	@$(MIGRATE) drop

.PHONY: dry-run
dry-run: migrate-drop run-migrate

.PHONY: gen-api
gen-api:
	rm -rf internal/gen/api/*
	protoc \
		-I api \
		-I api/validate \
		--go_out=internal/gen/api \
		--go_opt=paths=source_relative \
		--twirp_out=internal/gen/api \
		--twirp_opt=paths=source_relative \
		--validate_out="lang=go,paths=source_relative:internal/gen/api" \
		api/tag/v1/*.proto
	protoc \
		-I api \
		-I api/validate \
		--go_out=internal/gen/api \
		--go_opt=paths=source_relative \
		--twirp_out=internal/gen/api \
		--twirp_opt=paths=source_relative \
		--validate_out="lang=go,paths=source_relative:internal/gen/api" \
		api/player/v1/*.proto
	protoc \
		-I api \
		-I api/validate \
		--go_out=internal/gen/api \
		--go_opt=paths=source_relative \
		--twirp_out=internal/gen/api \
		--twirp_opt=paths=source_relative \
		--validate_out="lang=go,paths=source_relative:internal/gen/api" \
		api/question/v1/*.proto
	protoc \
		-I api \
		-I api/validate \
		--go_out=internal/gen/api \
		--go_opt=paths=source_relative \
		--twirp_out=internal/gen/api \
		--twirp_opt=paths=source_relative \
		--validate_out="lang=go,paths=source_relative:internal/gen/api" \
		api/package/v1/*.proto
	protoc \
		-I api \
		-I api/validate \
		--go_out=internal/gen/api \
		--go_opt=paths=source_relative \
		--twirp_out=internal/gen/api \
		--twirp_opt=paths=source_relative \
		--validate_out="lang=go,paths=source_relative:internal/gen/api" \
		api/auth/v1/*.proto