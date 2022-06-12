package psql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/answersuck/vault/internal/domain/account"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

type accountVerifRepo struct {
	l logging.Logger
	c *postgres.Client
}

func NewAccountVerificationRepo(l logging.Logger, c *postgres.Client) *accountVerifRepo {
	return &accountVerifRepo{
		l: l,
		c: c,
	}
}

func (r *accountVerifRepo) Verify(ctx context.Context, code string, updatedAt time.Time) error {
	sql := `
		UPDATE account a
		SET
			is_verified = $1,
			updated_at = $2
		FROM (
			SELECT v.code, a.id AS account_id
			FROM verification v
			INNER JOIN account a
			ON v.account_id = a.id
			WHERE v.code = $3
		) AS sq
		WHERE
			a.is_verified = $4
			AND a.id = sq.account_id;
	`

	r.l.Info("psql - account - Verify: %s", sql)

	ct, err := r.c.Pool.Exec(ctx, sql, true, updatedAt, code, false)
	if err != nil {
		return fmt.Errorf("psql - account_verification - Verify - r.c.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("psql - account_verification - Verify - r.c.Pool.Exec: %w", account.ErrAlreadyVerified)
	}

	return nil
}

func (r *accountVerifRepo) Find(ctx context.Context, accountId string) (account.Verification, error) {
	sql := `
		SELECT 
			a.email, 
			a.is_verified, 
			v.code AS verification_code
		FROM account a
		LEFT JOIN verification v
		ON v.account_id = a.id
		WHERE a.id = $1
	`

	r.l.Info("psql - account_verification - Find: %s", sql)

	var v account.Verification

	if err := r.c.Pool.QueryRow(ctx, sql, accountId).Scan(
		&v.Email,
		&v.Verified,
		&v.Code,
	); err != nil {

		if err == pgx.ErrNoRows {
			return account.Verification{},
				fmt.Errorf("psql - account_verification - Find - r.c.Pool.QueryRow.Scan: %w", account.ErrVerificationNotFound)
		}

		return account.Verification{}, fmt.Errorf("psql - account_verification - Find - r.c.Pool.QueryRow.Scan: %w", err)
	}

	return v, nil
}
