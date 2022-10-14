package migrate

import (
	"errors"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultAttempts = 10
	defaultTimeout  = time.Second
)

type operation string

const (
	Up   = operation("up")
	Down = operation("down")
)

func Do(op operation, dir, connString string) {
	m, err := connect(dir, connString)
	if err != nil {
		log.Fatalf("migrate: postgres connect error: %s", err.Error())
	}
	defer m.Close()

	switch op {
	case Up:
		err = m.Up()
	case Down:
		err = m.Down()
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrate: %s error: %s", string(op), err.Error())
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("migrate: %s no change", string(op))
	}
}

func connect(dir, connString string) (*migrate.Migrate, error) {
	connString += "?sslmode=disable"

	var (
		attempts = defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://"+dir, connString)
		if err != nil {
			return nil, err
		}
		if err == nil {
			break
		}

		log.Printf("migrate: trying connecting to postgres, attempts left: %d", attempts)
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		return nil, err
	}

	return m, nil
}
