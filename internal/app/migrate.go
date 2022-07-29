package app

import "github.com/answersuck/vault/pkg/migrate"

func init() { migrate.Up("file://migrations") }
