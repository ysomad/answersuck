package repository_psql

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"

	"github.com/ysomad/answersuck-backend/internal/adapter/repository/psql"
	"github.com/ysomad/answersuck-backend/internal/pkg/logger"
	"github.com/ysomad/answersuck-backend/internal/pkg/migrate"
	"github.com/ysomad/answersuck-backend/internal/pkg/postgres"
)

func initRepos(logLevel string) {
	log := logger.New(os.Stdout, logLevel)
	_accountRepo = psql.NewAccountRepo(log, db)
	_sessionRepo = psql.NewSessionRepo(log, db)
	_mediaRepo = psql.NewMediaRepo(log, db)
	_tagRepo = psql.NewTagRepo(log, db)
}

var db *postgres.Client

func TestMain(m *testing.M) {
	var err error

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14-alpine",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=test",
			"POSTGRES_DB=test",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://test:secret@%s/test?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseURL)

	resource.Expire(120)

	db, err = postgres.NewClient(databaseURL)
	if err != nil {
		log.Fatalf("Error initializing Postgres test client: %v", err)
	}
	defer db.Pool.Close()

	migrate.Up("../../migrations")

	initRepos(os.Getenv("LOG_LEVEL"))

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
