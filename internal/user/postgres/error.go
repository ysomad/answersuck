package postgres

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ysomad/answersuck/internal/user/domain"
)

const _errFmt = "%s: got '%s', returned '%w'"

// errUniqueViolation returns domain error wrapped with original error depending on constraint name,
// return original error if domain errors with given constraint name not found.
//
// Returned error format: {msg}: original error '{err}', returned error '{domain.Err}'.
func errUniqueViolation(msg string, err *pgconn.PgError) error {
	switch err.ConstraintName {
	case "account_email_key":
		return fmt.Errorf(_errFmt, msg, err.Error(), domain.ErrEmailTaken)
	case "account_username_key":
		return fmt.Errorf(_errFmt, msg, err.Error(), domain.ErrUsernameTaken)
	}

	return err
}
