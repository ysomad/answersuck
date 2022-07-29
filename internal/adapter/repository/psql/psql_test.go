package psql_test

import (
	"log"
	"os"
	"testing"

	"github.com/answersuck/vault/pkg/migrate"
)

func TestMain(m *testing.M) {
	if os.Getenv("INTEGRATION_TESTDB") != "true" {
		log.Printf("Skipping tests that require database connection")
		return
	}
	u := "file://../../../../migrations"
	migrate.Up(u)
	c := m.Run()
	migrate.Down(u)
	os.Exit(c)
}
