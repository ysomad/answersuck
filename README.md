# Answersuck backend API application

## Directories description

### `api/`

Directory for API documentation.

### `cmd/`

Entrypoint for the project. Can contain multiple entrypoints for microservices for example.

### `internal/`

Private application code.

### `pkg/`

Public application code that's ok to use by external applications.
Other projects will import these libraries expecting them to work, so think twice before you put something here :-)

### `migrations/`

Database migrations.

### `web/`

Web templates.

### `test/`

Directory for integration tests. Rest of tests are better to keep next to code which is needed to be tested.

## Local development
1. Configure `env` file
2. Run services from docker compose
```shell
$ make compose-up
```
3. Compile and run the app
```shell
$ make run
```

Make sure Makefile includes correct `env` file.

## Documentation
OpenAPI 3.0 documentation hosted in Swagger UI at `{host}:{port}/docs`
