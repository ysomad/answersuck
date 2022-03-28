package repository

import (
	"context"
	"fmt"
	"github.com/answersuck/answersuck-backend/internal/dto"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/answersuck/answersuck-backend/internal/domain"
	"github.com/answersuck/answersuck-backend/pkg/postgres"
)

const (
	accountTable             = "account"
	accountAvatarTable       = "account_avatar"
	accountVerificationTable = "account_verification"
)

type accountRepository struct {
	*postgres.Client
}

func NewAccountRepository(pg *postgres.Client) *accountRepository {
	return &accountRepository{pg}
}

func (r *accountRepository) Create(ctx context.Context, a *domain.Account) (*domain.Account, error) {
	sql := fmt.Sprintf(`
		WITH a AS (
			INSERT INTO %s(username, email, password, is_verified)
			VALUES ($1, $2, $3, $4)
			RETURNING id AS account_id
		),
		av AS (
			INSERT INTO %s(code, account_id)
			VALUES($5, (SELECT account_id FROM a) )
		),
		aa AS (
			INSERT INTO %s(url, account_id)
			VALUES($6, (SELECT account_id FROM a) )
		)
		SELECT account_id FROM a
	`, accountTable, accountVerificationTable, accountAvatarTable)

	if err := r.Pool.QueryRow(ctx, sql,
		a.Username,
		a.Email,
		a.PasswordHash,
		a.Verified,
		a.VerificationCode,
		a.AvatarURL,
	).Scan(&a.Id); err != nil {
		if err = isUniqueViolation(err); err != nil {
			return nil, fmt.Errorf("r.QueryRow.Exec: %w", err)
		}

		return nil, fmt.Errorf("r.QueryRow.Exec: %w", err)
	}

	return a, nil
}

func (r *accountRepository) FindById(ctx context.Context, aid string) (*domain.Account, error) {
	sql := fmt.Sprintf(`
		SELECT
			a.username,
			a.email,
			a.password,
			a.created_at,
			a.updated_at,
			a.is_verified,
			aa.url AS avatar_url
		FROM %s AS a
		LEFT JOIN %s AS aa
		ON aa.account_id = a.id
		WHERE
			a.id = $1
			AND a.is_archived = $2
	`, accountTable, accountAvatarTable)

	a := domain.Account{Id: aid}

	if err := r.Pool.QueryRow(ctx, sql, aid, false).Scan(
		&a.Username,
		&a.Email,
		&a.PasswordHash,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.Verified,
		&a.AvatarURL,
	); err != nil {

		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return &a, nil
}

func (r *accountRepository) FindByEmail(ctx context.Context, email string) (*domain.Account, error) {
	sql := fmt.Sprintf(`
		SELECT
			a.id,
			a.username,
			a.password,
			a.created_at,
			a.updated_at,
			a.is_verified,
			aa.url AS avatar_url
		FROM %s AS a
		LEFT JOIN %s AS aa
		ON aa.account_id = a.id
		WHERE
			a.email = $1
			AND a.is_archived = $2
	`, accountTable, accountAvatarTable)

	a := domain.Account{Email: email}

	if err := r.Pool.QueryRow(ctx, sql, email, false).Scan(
		&a.Id,
		&a.Username,
		&a.PasswordHash,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.Verified,
		&a.AvatarURL,
	); err != nil {

		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return &a, nil
}

func (r *accountRepository) FindByUsername(ctx context.Context, username string) (*domain.Account, error) {
	sql := fmt.Sprintf(`
		SELECT
			a.id,
			a.email,
			a.password,
			a.created_at,
			a.updated_at,
			a.is_verified,
			aa.url AS avatar_url
		FROM %s AS a
		LEFT JOIN %s AS aa
		ON aa.account_id = a.id
		WHERE
			a.username = $1
			AND a.is_archived = $2
	`, accountTable, accountAvatarTable)

	a := domain.Account{Username: username}

	if err := r.Pool.QueryRow(ctx, sql, username, false).Scan(
		&a.Id,
		&a.Email,
		&a.PasswordHash,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.Verified,
		&a.AvatarURL,
	); err != nil {

		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return &a, nil
}

func (r *accountRepository) Archive(ctx context.Context, aid string, archived bool, updatedAt time.Time) error {
	sql := fmt.Sprintf(`
		UPDATE %s
		SET
			is_archived = $1,
			updated_at = $2
		WHERE
			id = $3
			AND is_archived = $4
	`, accountTable)

	ct, err := r.Pool.Exec(ctx, sql, archived, updatedAt, aid, !archived)
	if err != nil {
		return fmt.Errorf("r.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("r.Pool.Exec: %w", ErrNoAffectedRows)
	}

	return nil
}

func (r *accountRepository) Verify(ctx context.Context, code string, verified bool, updatedAt time.Time) error {
	sql := fmt.Sprintf(`
		UPDATE %s AS a
		SET
			is_verified = $1,
			updated_at = $2
		FROM %s AS av
		WHERE
			a.is_verified = $3
			AND av.code = $4
	`, accountTable, accountVerificationTable)

	ct, err := r.Pool.Exec(ctx, sql, verified, updatedAt, !verified, code)
	if err != nil {
		return fmt.Errorf("r.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("r.Pool.Exec: %w", ErrNoAffectedRows)
	}

	return nil
}

func (r *accountRepository) FindVerification(ctx context.Context, aid string) (dto.AccountVerification, error) {
	sql := fmt.Sprintf(`
		SELECT a.email, a.is_verified, av.code AS verification_code
		FROM %s AS a
		LEFT JOIN %s AS av
		ON av.account_id = a.id
		WHERE a.id = $1
	`, accountTable, accountVerificationTable)

	var a dto.AccountVerification

	if err := r.Pool.QueryRow(ctx, sql, aid).Scan(&a.Email, &a.Verified, &a.Code); err != nil {
		if err == pgx.ErrNoRows {
			return dto.AccountVerification{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", ErrNotFound)
		}

		return dto.AccountVerification{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return a, nil
}
