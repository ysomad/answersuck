//go:build migrate
// +build migrate

package app

import (
	"os"

	"github.com/ysomad/answersuck/migrate"
)

func init() {
	migrate.Do(migrate.Up, "../user/postgres/migrations", os.Getenv("PG_URL"))
}
