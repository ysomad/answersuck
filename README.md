# Uniplay

## Local development

Swagger will be served at `localhost:8080/v1/docs`;
Grafana - `localhost:3000`;
Jaeger - `localhost:16686`.

### Run service deps and app separately
1. Run service deps
```sh
$ make compose-up
```
2. Run the application
```sh
$ make run-migrate
```

### Run service and deps in docker
- Run only required for development services
```sh
$ make compose-min
```

- Run all services
```sh
$ make compose-all
```