//go:build migrate
// +build migrate

package app

import (
	"os"

	"github.com/ysomad/answersuck/migrate"
)

func init() {
	migrate.Do(migrate.Up, "../peasant/postgres/migrations", os.Getenv("PG_URL"))
}
