package psql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/answersuck/host/internal/domain/account"
	"github.com/answersuck/host/internal/pkg/postgres"
)

type AccountRepo struct {
	*zap.Logger
	*postgres.Client
}

func NewAccountRepo(l *zap.Logger, c *postgres.Client) *AccountRepo {
	return &AccountRepo{l, c}
}

func (r *AccountRepo) Save(ctx context.Context, a account.Account, code string) (account.Account, error) {
	sql := `
WITH a AS (
    INSERT INTO account(email, nickname, password, is_verified, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id AS account_id
), v AS (
    INSERT INTO verification(code, account_id)
    VALUES($7, (SELECT account_id FROM a) )
), p AS (
    INSERT INTO player(account_id)
    VALUES((SELECT account_id FROM a))
)
SELECT account_id FROM a`

	r.Debug("psql - account - Save", zap.String("sql", sql), zap.Any("account", a))

	err := r.Pool.QueryRow(ctx, sql,
		a.Email,
		a.Nickname,
		a.Password,
		a.Verified,
		a.CreatedAt,
		a.UpdatedAt,
		code,
	).Scan(&a.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return account.Account{}, fmt.Errorf("psql - account - Save - r.Pool.QueryRow.Scan: %w", account.ErrAlreadyExist)
			}
		}
		return account.Account{}, fmt.Errorf("psql - account - Save - r.Pool.QueryRow.Scan: %w", err)
	}

	return a, nil
}

func (r *AccountRepo) FindById(ctx context.Context, accountId string) (account.Account, error) {
	sql := `
SELECT email, nickname, password, is_verified, created_at, updated_at
FROM account
WHERE id = $1 AND is_archived = $2`

	r.Debug("psql - account - FindById", zap.String("sql", sql), zap.String("accountId", accountId))

	var a account.Account

	if err := r.Pool.QueryRow(ctx, sql, accountId, false).Scan(
		&a.Email,
		&a.Nickname,
		&a.Password,
		&a.Verified,
		&a.CreatedAt,
		&a.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return account.Account{}, fmt.Errorf("psql - account - FindById - r.Pool.QueryRow.Scan: %w", account.ErrNotFound)
		}

		return account.Account{}, fmt.Errorf("psql - account - FindById - r.Pool.QueryRow.Scan: %w", err)
	}

	a.Id = accountId

	return a, nil
}

func (r *AccountRepo) FindByEmail(ctx context.Context, email string) (account.Account, error) {
	sql := `
SELECT id, nickname, password, created_at, updated_at, is_verified
FROM account
WHERE email = $1 AND is_archived = $2`

	r.Debug("psql - account - FindByEmail", zap.String("sql", sql), zap.String("email", email))

	var a account.Account

	if err := r.Pool.QueryRow(ctx, sql, email, false).Scan(
		&a.Id,
		&a.Nickname,
		&a.Password,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.Verified,
	); err != nil {

		if err == pgx.ErrNoRows {
			return account.Account{}, fmt.Errorf("psql - account - FindByEmail - r.Pool.QueryRow.Scan: %w", account.ErrNotFound)
		}

		return account.Account{}, fmt.Errorf("psql - account - FindByEmail - r.Pool.QueryRow.Scan: %w", err)
	}

	a.Email = email

	return a, nil
}

func (r *AccountRepo) FindByNickname(ctx context.Context, nickname string) (account.Account, error) {
	sql := `
SELECT id, email, password, created_at, updated_at, is_verified
FROM account
WHERE nickname = $1 AND is_archived = $2`

	r.Debug("psql - account - FindByNickname", zap.String("sql", sql), zap.String("nickname", nickname))

	var a account.Account

	err := r.Pool.QueryRow(ctx, sql, nickname, false).Scan(
		&a.Id,
		&a.Email,
		&a.Password,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.Verified,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return account.Account{}, fmt.Errorf("psql - account - FindByNickname - r.Pool.QueryRow.Scan: %w", account.ErrNotFound)
		}

		return account.Account{}, fmt.Errorf("psql - account - FindByNickname - r.Pool.QueryRow.Scan: %w", err)
	}

	a.Nickname = nickname

	return a, nil
}

func (r *AccountRepo) SetArchived(ctx context.Context, accountId string, archived bool, updatedAt time.Time) error {
	sql := `
UPDATE account SET is_archived = $1, updated_at = $2
WHERE id = $3 AND is_archived = $4`

	r.Debug("psql - account - SetArchived", zap.String("sql", sql), zap.String("accountId", accountId))

	ct, err := r.Pool.Exec(ctx, sql, archived, updatedAt, accountId, !archived)
	if err != nil {
		return fmt.Errorf("psql - account - SetArchived - r.Pool.Exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("psql - account - SetArchived - r.Pool.Exec: %w", account.ErrNotFound)
	}

	return nil
}

func (r *AccountRepo) FindPasswordById(ctx context.Context, accountId string) (string, error) {
	sql := "SELECT password FROM account WHERE id = $1"

	r.Debug("psql - account - FindPasswordById", zap.String("accountId", accountId))

	var password string
	if err := r.Pool.QueryRow(ctx, sql, accountId).Scan(&password); err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("psql - account - FindPasswordById - r.Pool.QueryRow.Scan: %w", account.ErrNotFound)
		}

		return "", fmt.Errorf("psql - account - FindPasswordById - r.Pool.QueryRow.Scan: %w", err)
	}

	return password, nil
}

func (r *AccountRepo) UpdatePassword(ctx context.Context, accountId, password string) error {
	sql := "UPDATE account SET password = $1 WHERE id = $2"

	r.Debug(
		"psql - account - UpdatePassword",
		zap.String("accountId", accountId),
		zap.String("password", password),
	)

	ct, err := r.Pool.Exec(ctx, sql, password, accountId)
	if err != nil {
		return fmt.Errorf("psql - account - UpdatePassword - r.Pool.Exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("psql - account - UpdatePassword - r.Pool.Exec: %w", account.ErrNotFound)
	}

	return nil
}

func (r *AccountRepo) SavePasswordToken(ctx context.Context, dto account.SavePasswordTokenDTO) (string, error) {
	sql := `
WITH e AS (
    SELECT id AS account_id, email AS account_email FROM account
    WHERE email = $1 OR nickname = $2
), pt AS (
    INSERT INTO password_token(token, account_id)
    VALUES($3, (SELECT account_id FROM e))
)
SELECT account_email FROM e`

	r.Debug(
		"psql - account - SavePasswordToken",
		zap.String("sql", sql),
		zap.String("login", dto.Login),
		zap.String("token", dto.Token),
		zap.Time("created_at", dto.CreatedAt),
	)

	var email string
	if err := r.Pool.QueryRow(ctx, sql, dto.Login, dto.Login, dto.Token).Scan(&email); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return "",
					fmt.Errorf("psql - account - SavePasswordToken - r.Pool.QueryRow.Scan: %w", account.ErrPasswordTokenAlreadyExist)
			case pgerrcode.NotNullViolation:
				return "",
					fmt.Errorf("psql - account - SavePasswordToken - r.Pool.QueryRow.Scan: %w", account.ErrNotFound)
			}
		}

		return "", fmt.Errorf("psql - account - SavePasswordToken - r.QueryRow.Scan: %w", err)
	}

	return email, nil
}

func (r *AccountRepo) FindPasswordToken(ctx context.Context, token string) (account.PasswordToken, error) {
	sql := `
SELECT t.token, t.created_at, a.id
FROM password_token t
INNER JOIN account a
ON a.id = t.account_id
WHERE t.token = $1`

	r.Debug("psql - account - FindPasswordToken", zap.String("sql", sql), zap.String("token", token))

	var t account.PasswordToken

	if err := r.Pool.QueryRow(ctx, sql, token).Scan(&t.Token, &t.CreatedAt, &t.AccountId); err != nil {
		if err == pgx.ErrNoRows {
			return account.PasswordToken{},
				fmt.Errorf("psql - account - FindPasswordToken - r.Pool.QueryRow.Scan: %w", account.ErrPasswordTokenNotFound)
		}

		return account.PasswordToken{},
			fmt.Errorf("psql - account - FindPasswordToken - r.client.Pool.QueryRow.Scan: %w", err)
	}

	return t, nil
}

func (r *AccountRepo) SetPassword(ctx context.Context, dto account.SetPasswordDTO) error {
	sql := `
WITH a AS (
    UPDATE account
    SET password = $1, updated_at = $2
    WHERE id = $3
)
DELETE FROM password_token
WHERE account_id = $4 AND token = $5`

	r.Debug("psql - account - SetPassword", zap.String("sql", sql), zap.Any("args", dto))

	ct, err := r.Pool.Exec(
		ctx,
		sql,
		dto.Password,
		dto.UpdatedAt,
		dto.AccountId,
		dto.AccountId,
		dto.Token,
	)
	if err != nil {
		return fmt.Errorf("psql - account - SetPassword - r.Pool.Exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("psql - account - SetPassword - r.Pool.Exec: %w", account.ErrPasswordNotSet)
	}

	return nil
}

func (r *AccountRepo) Verify(ctx context.Context, code string, updatedAt time.Time) error {
	sql := `
UPDATE account a
SET is_verified = $1, updated_at = $2
FROM (
    SELECT v.code, a.id AS account_id
    FROM verification v
    INNER JOIN account a
    ON v.account_id = a.id
    WHERE v.code = $3
) AS sq
WHERE a.is_verified = $4 AND a.id = sq.account_id`

	r.Debug("psql - account - Verify", zap.String("sql", sql), zap.String("code", code))

	ct, err := r.Pool.Exec(ctx, sql, true, updatedAt, code, false)
	if err != nil {
		return fmt.Errorf("psql - account - Verify - r.Pool.Exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("psql - account - Verify - r.Pool.Exec: %w", account.ErrAlreadyVerified)
	}

	return nil
}

func (r *AccountRepo) FindVerification(ctx context.Context, accountId string) (account.Verification, error) {
	sql := `
SELECT a.email, a.is_verified, v.code AS verification_code
FROM account a
LEFT JOIN verification v
ON v.account_id = a.id
WHERE a.id = $1`

	r.Debug("psql - account - FindVerification", zap.String("sql", sql), zap.String("accountId", accountId))

	var v account.Verification
	if err := r.Pool.QueryRow(ctx, sql, accountId).Scan(
		&v.Email,
		&v.Verified,
		&v.Code,
	); err != nil {
		if err == pgx.ErrNoRows {
			return account.Verification{},
				fmt.Errorf("psql - account - FindVerification - r.Pool.QueryRow.Scan: %w", account.ErrVerificationNotFound)
		}

		return account.Verification{}, fmt.Errorf("psql - account - FindVerification - r.Pool.QueryRow.Scan: %w", err)
	}

	return v, nil
}
