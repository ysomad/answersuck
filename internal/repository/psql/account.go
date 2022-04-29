package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/internal/dto"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

const (
	accountTable                   = "account"
	accountAvatarTable             = "account_avatar"
	accountVerificationCodeTable   = "account_verification_code"
	accountPasswordResetTokenTable = "account_password_reset_token"
)

type account struct {
	log    logging.Logger
	client *postgres.Client
}

func NewAccount(l logging.Logger, c *postgres.Client) *account {
	return &account{
		log:    l,
		client: c,
	}
}

func (r *account) Create(ctx context.Context, a *domain.Account) (*domain.Account, error) {
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
	`, accountTable, accountVerificationCodeTable, accountAvatarTable)

	r.log.Info("psql - account - Create: %s", sql)

	if err := r.client.Pool.QueryRow(ctx, sql,
		a.Username,
		a.Email,
		a.PasswordHash,
		a.Verified,
		a.VerificationCode,
		a.AvatarURL,
	).Scan(&a.Id); err != nil {
		if err = unwrapError(err); err != nil {
			return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
		}

		return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return a, nil
}

func (r *account) FindById(ctx context.Context, aid string) (*domain.Account, error) {
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

	r.log.Info("psql - account - FindById: %s", sql)

	a := domain.Account{Id: aid}

	if err := r.client.Pool.QueryRow(ctx, sql, aid, false).Scan(
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

func (r *account) FindByEmail(ctx context.Context, email string) (*domain.Account, error) {
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

	r.log.Info("psql - account - FindByEmail: %s", sql)

	a := domain.Account{Email: email}

	if err := r.client.Pool.QueryRow(ctx, sql, email, false).Scan(
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

func (r *account) FindByUsername(ctx context.Context, username string) (*domain.Account, error) {
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

	r.log.Info("psql - account - FindByUsername: %s", sql)

	a := domain.Account{Username: username}

	if err := r.client.Pool.QueryRow(ctx, sql, username, false).Scan(
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

func (r *account) Archive(ctx context.Context, aid string, archived bool, updatedAt time.Time) error {
	sql := fmt.Sprintf(`
		UPDATE %s
		SET
			is_archived = $1,
			updated_at = $2
		WHERE
			id = $3
			AND is_archived = $4
	`, accountTable)

	r.log.Info("psql - account - Archive: %s", sql)

	ct, err := r.client.Pool.Exec(ctx, sql, archived, updatedAt, aid, !archived)
	if err != nil {
		return fmt.Errorf("r.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("r.Pool.Exec: %w", ErrNoAffectedRows)
	}

	return nil
}

func (r *account) Verify(ctx context.Context, code string, verified bool, updatedAt time.Time) error {
	sql := fmt.Sprintf(`
		UPDATE %s AS a
		SET 
			is_verified = $1,
			updated_at = $2
		FROM (
			SELECT av.code, a.id AS account_id
			FROM %s AS av
			INNER JOIN %s AS a
			ON av.account_id = a.id
			WHERE av.code = $3
		) AS sq
		WHERE
			a.is_verified = $4 
			AND a.id = sq.account_id;
	`, accountTable, accountVerificationCodeTable, accountTable)

	r.log.Info("psql - account - Verify: %s", sql)

	ct, err := r.client.Pool.Exec(ctx, sql, verified, updatedAt, code, !verified)
	if err != nil {
		return fmt.Errorf("r.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("r.Pool.Exec: %w", ErrNoAffectedRows)
	}

	return nil
}

func (r *account) FindVerification(ctx context.Context, aid string) (dto.AccountVerification, error) {
	sql := fmt.Sprintf(`
		SELECT a.email, a.is_verified, av.code AS verification_code
		FROM %s AS a
		LEFT JOIN %s AS av
		ON av.account_id = a.id
		WHERE a.id = $1
	`, accountTable, accountVerificationCodeTable)

	r.log.Info("psql - account - FindVerification: %s", sql)

	var a dto.AccountVerification

	if err := r.client.Pool.QueryRow(ctx, sql, aid).Scan(&a.Email, &a.Verified, &a.Code); err != nil {
		if err == pgx.ErrNoRows {
			return dto.AccountVerification{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", ErrNotFound)
		}

		return dto.AccountVerification{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return a, nil
}

func (r *account) InsertPasswordResetToken(ctx context.Context, email, token string) error {
	sql := fmt.Sprintf(`
		INSERT INTO %s (token, account_id)
		VALUES($1, (SELECT id AS account_id FROM %s WHERE email = $2))
	`, accountPasswordResetTokenTable, accountTable)

	r.log.Info("psql - account - InsertPasswordResetToken: %s", sql)

	if _, err := r.client.Pool.Exec(ctx, sql, token, email); err != nil {

		if err = unwrapError(err); err != nil {
			return fmt.Errorf("r.Pool.Exec: %w", err)
		}

		return fmt.Errorf("r.Pool.Exec: %w", ErrNotFound)
	}

	return nil
}

func (r *account) FindPasswordResetToken(ctx context.Context, token string) (*dto.AccountPasswordResetToken, error) {
	sql := fmt.Sprintf(`
		SELECT t.token, t.created_at, a.id
		FROM %s AS t
		INNER JOIN %s AS a
		ON a.id = t.account_id
		WHERE t.token = $1
	`, accountPasswordResetTokenTable, accountTable)

	r.log.Info("psql - account - FindPasswordResetToken: %s", sql)

	var a dto.AccountPasswordResetToken

	if err := r.client.Pool.QueryRow(ctx, sql, token).Scan(&a.Token, &a.CreatedAt, &a.AccountId); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return &a, nil
}

func (r *account) UpdatePasswordWithToken(ctx context.Context, dto dto.AccountUpdatePassword) error {
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

	r.log.Info("psql - account - UpdatePasswordWithToken: %s", sql)

	ct, err := r.client.Pool.Exec(
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
