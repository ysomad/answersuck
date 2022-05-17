# Answersuck backend API application

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