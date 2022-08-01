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

var postgresClient *postgres.Client

func initRepos(logLevel string) {
	logger := logger.New(os.Stdout, logLevel)
	accountRepo = psql.NewAccountRepo(logger, postgresClient)
	sessionRepo = psql.NewSessionRepo(logger, postgresClient)
}

func TestMain(m *testing.M) {
	u := "../../migrations"
	migrate.Down(u)
	migrate.Up(u)

	postgresURI := os.Getenv("PG_URL")
	if postgresURI == "" {
		log.Fatal("Empty PG_URL environment variable")
	}

	var err error
	postgresClient, err = postgres.NewClient(postgresURI)
	if err != nil {
		log.Fatalf("Error initializing Postgres test client: %v", err)
	}

	initRepos(os.Getenv("LOG_LEVEL"))

	os.Exit(m.Run())
}
