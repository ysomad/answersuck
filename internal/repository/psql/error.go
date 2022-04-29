package repository

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

var (
	ErrUniqueViolation = errors.New("duplicate key value violates unique constraint")
	ErrNotFound        = errors.New("requested entity not found in database")
	ErrNoAffectedRows  = errors.New("zero rows affected")
)

func unwrapError(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return ErrUniqueViolation
		}
	}

	return nil
}
