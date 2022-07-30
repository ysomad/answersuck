package psql_test

import (
	"log"
	"os"
	"testing"

	"github.com/answersuck/vault/internal/adapter/repository/psql"

	"github.com/answersuck/vault/pkg/logger"
	"github.com/answersuck/vault/pkg/migrate"
	"github.com/answersuck/vault/pkg/postgres"
)

var postgresClient *postgres.Client

func initRepos(logLevel string) {
	logger := logger.New(os.Stdout, logLevel)
	accountRepo = psql.NewAccountRepo(logger, postgresClient)
}

func TestMain(m *testing.M) {
	if os.Getenv("INTEGRATION_TESTDB") != "true" {
		log.Printf("Skipping tests that require database connection")
		return
	}

	u := "../../../../migrations"
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

	initRepos(os.Getenv("INTEGRATION_LOGLEVEL"))

	c := m.Run()
	os.Exit(c)
}
