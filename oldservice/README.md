# Answersuck backend API application

## Local development

1. Configure `.env` file
2. Run services from docker compose
```shell
$ make compose-up
```
3. Compile and run the app
```shell
$ make run
```

## Run tests locally

- Run unit tests `make test`
- Run integration tests `make integration-test`

Make sure required services for integration tests are running in docker-compose.

## Directories description

### `api/`

Directory for API documentation.

### `cmd/`

Entrypoint for the project. Can contain multiple entrypoints for microservices for example.

### `internal/`

Private application code.

### `migrations/`

Database migrations.

### `web/`

Web templates.

### `test/`

Directory for integration tests. Rest of tests are better to keep next to code which is needed to be tested.
