package repository_psql

import (
	"log"
	"os"
	"testing"

	"github.com/answersuck/host/internal/adapter/repository/psql"
	"github.com/answersuck/host/internal/pkg/logger"
	"github.com/answersuck/host/internal/pkg/migrate"
	"github.com/answersuck/host/internal/pkg/postgres"
)

func initRepos(logLevel string, c *postgres.Client) {
	logger := logger.New(os.Stdout, logLevel)
	_accountRepo = psql.NewAccountRepo(logger, c)
	_sessionRepo = psql.NewSessionRepo(logger, c)
	_mediaRepo = psql.NewMediaRepo(logger, c)
}

func TestMain(m *testing.M) {
	u := "../../migrations"
	migrate.Down(u)
	migrate.Up(u)

	postgresURI := os.Getenv("PG_URL")
	if postgresURI == "" {
		log.Fatal("Empty PG_URL environment variable")
	}

	postgresClient, err := postgres.NewClient(postgresURI)
	if err != nil {
		log.Fatalf("Error initializing Postgres test client: %v", err)
	}

	initRepos(os.Getenv("LOG_LEVEL"), postgresClient)

	os.Exit(m.Run())
}
