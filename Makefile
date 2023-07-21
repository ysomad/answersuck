include .env
export

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=${PG_URL}

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
	go run ./cmd/app -migrate

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -v -cover -race -count 1 ./internal/...

.PHONY: goose-new
goose-new:
	@read -p "Enter the name of the new migration: " name; \
	goose -s -dir migrations create $${name// /_} sql

.PHONY: goose-up
goose-up:
	@echo "Running all new database migrations..."
	goose -dir migrations validate
	goose -dir migrations up

.PHONY: goose-down
goose-down:
	@echo "Running all down database migrations..."
	goose -dir migrations down

.PHONY: goose-reset
goose-reset:
	@echo "Dropping everything in database..."
	goose -dir migrations reset

.PHONY: goose-status
goose-status:
	goose -dir migrations status

.PHONY: dry-run
dry-run: goose-reset run-migrate

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