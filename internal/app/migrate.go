//go:build migrate
// +build migrate

package app

import (
	"log"
	"os"

	"github.com/ysomad/answersuck/internal/pkg/migrate"
)

func init() {
	connStr := os.Getenv("PG_URL")
	if connStr == "" {
		log.Fatal("PG_URL environment variable is not declared")
	}

	migrate.Do(migrate.Up, "./migrations", connStr)
}
