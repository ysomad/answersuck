package psql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	"github.com/answersuck/vault/internal/domain/account"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

type accountRepo struct {
	l logging.Logger
	c *postgres.Client
}

func NewAccountRepo(l logging.Logger, c *postgres.Client) *accountRepo {
	return &accountRepo{
		l: l,
		c: c,
	}
}

func (r *accountRepo) Save(ctx context.Context, a *account.Account, code string) (string, error) {
	sql := `
		WITH a AS (
			INSERT INTO account(email, nickname, password, is_verified, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id AS account_id
		),
		v AS (
			INSERT INTO verification(code, account_id)
			VALUES($7, (SELECT account_id FROM a) )
		),
		p AS (
			INSERT INTO player(account_id)
			VALUES((SELECT account_id FROM a))
		)
		SELECT account_id FROM a
	`

	r.l.Info("psql - account - Save: %s", sql)

	var accountId string

	err := r.c.Pool.QueryRow(ctx, sql,
		a.Email,
		a.Nickname,
		a.Password,
		a.Verified,
		a.CreatedAt,
		a.UpdatedAt,
		code,
	).Scan(&accountId)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			if pgErr.Code == pgerrcode.UniqueViolation {
				return "", fmt.Errorf("psql - account - Save - r.c.Pool.QueryRow.Scan: %w", account.ErrAlreadyExist)
			}

		}

		return "", fmt.Errorf("psql - account - Save - r.c.Pool.QueryRow.Scan: %w", err)
	}

	return accountId, nil
}

func (r *accountRepo) FindById(ctx context.Context, accountId string) (*account.Account, error) {
	sql := `
		SELECT
			email,
			nickname,
			password,
			is_verified,
			created_at,
			updated_at
		FROM account
		WHERE
			id = $1
			AND is_archived = $2
	`

	r.l.Info("psql - account - FindById: %s", sql)

	var a account.Account

	if err := r.c.Pool.QueryRow(ctx, sql, accountId, false).Scan(
		&a.Email,
		&a.Nickname,
		&a.Password,
		&a.Verified,
		&a.CreatedAt,
		&a.UpdatedAt,
	); err != nil {

		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("psql - account - FindById - r.c.Pool.QueryRow.Scan: %w", account.ErrNotFound)
		}

		return nil, fmt.Errorf("psql - account - FindById - r.c.Pool.QueryRow.Scan: %w", err)
	}

	a.Id = accountId

	return &a, nil
}

func (r *accountRepo) FindByEmail(ctx context.Context, email string) (*account.Account, error) {
	sql := `
		SELECT
			id,
			nickname,
			password,
			created_at,
			updated_at,
			is_verified
		FROM account
		WHERE
			email = $1
			AND is_archived = $2
	`

	r.l.Info("psql - account - FindByEmail: %s", sql)

	var a account.Account

	if err := r.c.Pool.QueryRow(ctx, sql, email, false).Scan(
		&a.Id,
		&a.Nickname,
		&a.Password,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.Verified,
	); err != nil {

		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("psql - account - FindByEmail - r.c.Pool.QueryRow.Scan: %w", account.ErrNotFound)
		}

		return nil, fmt.Errorf("psql - account - FindByEmail - r.client.Pool.QueryRow.Scan: %w", err)
	}

	a.Email = email

	return &a, nil
}

func (r *accountRepo) FindByNickname(ctx context.Context, nickname string) (*account.Account, error) {
	sql := `
		SELECT
			id,
			email,
			password,
			created_at,
			updated_at,
			is_verified
		FROM account
		WHERE
			nickname = $1
			AND is_archived = $2
	`

	r.l.Info("psql - account - FindByNickname: %s", sql)

	var a account.Account

	err := r.c.Pool.QueryRow(ctx, sql, nickname, false).Scan(
		&a.Id,
		&a.Email,
		&a.Password,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.Verified,
	)
	if err != nil {

		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("psql - account - FindByNickname - r.c.Pool.QueryRow.Scan: %w", account.ErrNotFound)
		}

		return nil, fmt.Errorf("psql - account - FindByNickname - r.c.Pool.QueryRow.Scan: %w", err)
	}

	a.Nickname = nickname

	return &a, nil
}

func (r *accountRepo) Archive(ctx context.Context, accountId string, updatedAt time.Time) error {
	sql := `
		UPDATE account
		SET
			is_archived = $1,
			updated_at = $2
		WHERE
			id = $3
			AND is_archived = $4
	`

	r.l.Info("psql - account - Archive: %s", sql)

	ct, err := r.c.Pool.Exec(ctx, sql, true, updatedAt, accountId, false)

	if err != nil {
		return fmt.Errorf("psql - account - Archive - r.c.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("psql - account - Archive - r.c.Pool.Exec: %w", account.ErrNotDeleted)
	}

	return nil
}

func (r *accountRepo) SavePasswordToken(ctx context.Context, email, token string) error {
	sql := `
		INSERT INTO password_token(token, account_id)
		VALUES($1, (SELECT id AS account_id FROM account WHERE email = $2))
	`

	r.l.Info("psql - account - SavePasswordToken: %s", sql)

	if _, err := r.c.Pool.Exec(ctx, sql, token, email); err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return fmt.Errorf("psql - account - SavePasswordToken - r.c.Pool.Exec: %w", account.ErrPasswordTokenAlreadyExist)
			case pgerrcode.ForeignKeyViolation:
				return fmt.Errorf("psql - account - SavePasswordToken - r.c.Pool.Exec: %w", account.ErrNotFound)
			}
		}

		return fmt.Errorf("psql - account - SavePasswordToken - r.c.Pool.Exec: %w", err)
	}

	return nil
}

func (r *accountRepo) FindPasswordToken(ctx context.Context, token string) (account.PasswordToken, error) {
	sql := `
		SELECT t.token, t.created_at, a.id
		FROM password_token t
		INNER JOIN account a
		ON a.id = t.account_id
		WHERE t.token = $1
	`

	r.l.Info("psql - account - FindPasswordToken: %s", sql)

	var t account.PasswordToken

	if err := r.c.Pool.QueryRow(ctx, sql, token).Scan(&t.Token, &t.CreatedAt, &t.AccountId); err != nil {

		if err == pgx.ErrNoRows {
			return account.PasswordToken{},
				fmt.Errorf("psql - account - FindPasswordToken - r.c.Pool.QueryRow.Scan: %w", account.ErrPasswordTokenNotFound)
		}

		return account.PasswordToken{}, fmt.Errorf("psql - account - FindPasswordToken - r.client.Pool.QueryRow.Scan: %w", err)
	}

	return t, nil
}

func (r *accountRepo) FindEmailByNickname(ctx context.Context, nickname string) (string, error) {
	sql := "SELECT email FROM account WHERE nickname = $1"

	r.l.Info("psql - account - FindEmailByNickname: %s", sql)

	var email string

	if err := r.c.Pool.QueryRow(ctx, sql, nickname).Scan(&email); err != nil {

		if err == pgx.ErrNoRows {
			return "",
				fmt.Errorf("psql - account - FindEmailByNickname - r.c.Pool.QueryRow.Scan: %w", account.ErrNotFound)
		}

		return "", fmt.Errorf("psql - account - FindPasswordToken - r.client.Pool.QueryRow.Scan: %w", err)
	}

	return email, nil
}

func (r *accountRepo) SetPassword(ctx context.Context, dto account.SetPasswordDTO) error {
	sql := `
		WITH a AS (
			UPDATE account
			SET 
				password = $1, 
				updated_at = $2
			WHERE id = $3
		)
		DELETE FROM password_token
		WHERE
			account_id = $4
			AND token = $5
	`

	r.l.Info("psql - account - SetPassword: %s", sql)

	ct, err := r.c.Pool.Exec(
		ctx,
		sql,
		dto.Password,
		dto.UpdatedAt,
		dto.AccountId,
		dto.AccountId,
		dto.Token,
	)
	if err != nil {
		return fmt.Errorf("psql - account - SetPassword - r.c.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("psql - account - SetPassword - r.c.Pool.Exec: %w", account.ErrPasswordNotSet)
	}

	return nil
}
