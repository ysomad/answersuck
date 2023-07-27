# Answersuck

SIGame but better!

## Local development

1. Run service deps
```sh
$ make compose-dev
```
2. Run the application
```sh
$ make run
```

Or:
- ```make dry-run``` to run application with clean database
- ```make run-migrate``` to run app and up migrations

## Generate handlers from protobuf
```sh
$ make gen-api
```

## Generate swagger from protobuf
```sh
$ make gen-swagger
```