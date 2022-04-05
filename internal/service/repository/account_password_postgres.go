package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/answersuck/vault/internal/dto"

	"github.com/answersuck/vault/pkg/postgres"
)

const (
	accountPasswordResetTokenTable = "account_password_reset_token"
)

type accountPasswordRepository struct {
	*postgres.Client
}

func NewAccountPasswordRepository(pg *postgres.Client) *accountPasswordRepository {
	return &accountPasswordRepository{pg}
}

func (r *accountPasswordRepository) InsertResetToken(ctx context.Context, email, token string) error {
	sql := fmt.Sprintf(`
		INSERT INTO %s (token, account_id)
		VALUES($1, (SELECT id AS account_id FROM %s WHERE email = $2))
	`, accountPasswordResetTokenTable, accountTable)

	if _, err := r.Pool.Exec(ctx, sql, token, email); err != nil {

		if err = isUniqueViolation(err); err != nil {
			return fmt.Errorf("r.Pool.Exec: %w", err)
		}

		return fmt.Errorf("r.Pool.Exec: %w", ErrNotFound)
	}

	return nil
}

func (r *accountPasswordRepository) FindResetToken(ctx context.Context, token string) (*dto.AccountPasswordResetToken, error) {
	sql := fmt.Sprintf(`
		SELECT t.token, t.created_at, a.id
		FROM %s AS t
		INNER JOIN %s AS a
		ON a.id = t.account_id
		WHERE t.token = $1
	`, accountPasswordResetTokenTable, accountTable)

	var a dto.AccountPasswordResetToken

	if err := r.Pool.QueryRow(ctx, sql, token).Scan(&a.Token, &a.CreatedAt, &a.AccountId); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return &a, nil
}

func (r *accountPasswordRepository) UpdateWithToken(ctx context.Context, dto dto.AccountUpdatePassword) error {
	sql := fmt.Sprintf(`
		WITH a AS (
			UPDATE %s
			SET password = $1, updated_at = $2
			WHERE id = $3
		)
		DELETE FROM %s
		WHERE
			account_id = $4
			AND token = $5
	`, accountTable, accountPasswordResetTokenTable)

	ct, err := r.Pool.Exec(
		ctx,
		sql,
		dto.PasswordHash,
		dto.UpdatedAt,
		dto.AccountId,
		dto.AccountId,
		dto.Token,
	)
	if err != nil {
		return fmt.Errorf("r.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("r.Pool.Exec: %w", ErrNoAffectedRows)
	}

	return nil
}
