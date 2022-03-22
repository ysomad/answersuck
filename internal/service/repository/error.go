package repository

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

var (
	ErrUniqueViolation = errors.New("unique violation error")
	ErrNotFound        = errors.New("requested entity not found in database")
)

func isUniqueViolation(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return ErrUniqueViolation
		}
	}

	return nil
}
