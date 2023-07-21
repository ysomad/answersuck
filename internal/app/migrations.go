package app

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func mustRunMigrations(dbURL string) {
	db, err := goose.OpenDBWithDriver("pgx", dbURL)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	if err := goose.Run("up", db, "./migrations"); err != nil {
		panic(err)
	}
}
