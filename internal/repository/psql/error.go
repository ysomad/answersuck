package repository

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

var (
	ErrUniqueViolation     = errors.New("duplicate key value violates unique constraint")
	ErrForeignKeyViolation = errors.New("entity with given foreign key is not present in database")
	ErrNotFound            = errors.New("requested entity not found in database")
	ErrNoAffectedRows      = errors.New("zero rows affected")
)

func unwrapError(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return ErrUniqueViolation
		case pgerrcode.ForeignKeyViolation:
			return ErrForeignKeyViolation
		}
	}

	return nil
}
