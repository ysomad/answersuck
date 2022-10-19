package postgres

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/ysomad/answersuck/apperror"
	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"

	"github.com/ysomad/answersuck/pgclient"
)

const accountReturnAll = "RETURNING id, username, email, is_email_verified, is_archived, created_at, updated_at"

type accountRepository struct {
	*pgclient.Client
}

func NewAccountRepository(c *pgclient.Client) *accountRepository {
	return &accountRepository{c}
}

func (r *accountRepository) Create(ctx context.Context, accArgs dto.AccountCreateArgs, verifArgs dto.EmailVerifCreateArgs) (*domain.Account, error) {
	const errMsg = "accountRepository - Create"

	insertAccountSQL, insertAccountArgs, err := r.Builder.
		Insert("account").
		Columns("email, username, password").
		Values(accArgs.Email, accArgs.Username, accArgs.Password).
		Suffix(accountReturnAll).
		ToSql()
	if err != nil {
		return nil, err
	}

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, insertAccountSQL, insertAccountArgs...)
	if err != nil {
		return nil, err
	}

	account, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[domain.Account])
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			switch pgErr.ConstraintName {
			case constraintAccountEmail:
				return nil, apperror.New(errMsg, err, domain.ErrEmailTaken)
			case constraintAccountUsername:
				return nil, apperror.New(errMsg, err, domain.ErrUsernameTaken)
			}
		}

		return nil, err
	}

	insertEmailVerifSQL, insertEmailVerifArgs, err := r.Builder.
		Insert("email_verification").
		Columns("account_id, code, expires_at").
		Values(account.ID, verifArgs.Code, verifArgs.ExpiresAt).
		ToSql()
	if err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, insertEmailVerifSQL, insertEmailVerifArgs...); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *accountRepository) GetByID(ctx context.Context, accountID string) (*domain.Account, error) {
	const errMsg = "accountRepository - GetByID"

	sql, args, err := r.Builder.
		Select("id, username, email, is_email_verified, is_archived, created_at, updated_at").
		From("account").
		Where(sq.Eq{"id": accountID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	a, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[domain.Account])
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.New(errMsg, err, domain.ErrAccountNotFound)
		}

		return nil, err
	}

	return &a, nil
}

func (r *accountRepository) DeleteByID(ctx context.Context, accountID string) error {
	sql, args, err := r.Builder.
		Update("account").
		Set("is_archived", true).
		Where(sq.And{
			sq.Eq{"id": accountID},
			sq.Eq{"is_archived": false},
		}).
		ToSql()
	if err != nil {
		return err
	}

	ct, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return domain.ErrAccountNotFound
	}

	return nil
}

func (r *accountRepository) GetPasswordByID(ctx context.Context, accountID string) (string, error) {
	const errMsg = "accountRepository - GetPasswordByID"

	sql, args, err := r.Builder.
		Select("password").
		From("account").
		Where(sq.Eq{"id": accountID}).
		ToSql()
	if err != nil {
		return "", err
	}

	var p string

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&p); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", apperror.New(errMsg, err, domain.ErrAccountNotFound)
		}

		return "", err
	}

	return p, nil
}

func (r *accountRepository) UpdateEmail(ctx context.Context, accountID, newEmail string) (*domain.Account, error) {
	const errMsg = "accountRepository - UpdateEmail"

	sql, args, err := r.Builder.
		Update("account").
		Set("email", newEmail).
		Set("is_email_verified", false).
		Where(sq.And{
			sq.Eq{"id": accountID},
			sq.Eq{"is_archived": false},
		}).
		Suffix("RETURNING username, email, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	// TODO: rewrite it using CollectOneRow when issue is fixed
	// https://github.com/jackc/pgx/issues/1334
	// rows, err := r.Pool.Query(ctx, query, args...)
	// if err != nil {
	// 	return nil, err
	// }

	// account, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[domain.Account])
	// if err != nil {
	// 	if errors.Is(err, pgx.ErrNoRows) {
	// 		return nil, apperror.New(errMsg, err, domain.ErrAccountNotFound)
	// 	}

	// 	var pgErr *pgconn.PgError
	// 	if errors.As(err, &pgErr) && pgErr.ConstraintName == constraintAccountEmail {
	// 		return nil, apperror.New(errMsg, pgErr, domain.ErrEmailTaken)
	// 	}

	// 	return nil, err
	// }

	var a domain.Account

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&a.Username,
		&a.Email,
		&a.CreatedAt,
		&a.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.New(errMsg, err, domain.ErrAccountNotFound)
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == constraintAccountEmail {
			return nil, apperror.New(errMsg, pgErr, domain.ErrEmailTaken)
		}

		return nil, err
	}

	a.ID, a.Email = accountID, newEmail
	return &a, nil
}

func (r *accountRepository) VerifyEmail(ctx context.Context, verifCode string) (*domain.Account, error) {
	const errMsg = "accountRepository - VerifyEmail"

	sql := `
UPDATE account SET is_email_verified = true
WHERE EXISTS(
(SELECT 1 FROM email_verification v
WHERE v.code = $1
AND v.account_id = account.id
AND v.expires_at < now())
) AND account.is_email_verified = false`

	rows, err := r.Pool.Query(ctx, sql, verifCode)
	if err != nil {
		return nil, err
	}

	a, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[domain.Account])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.New(errMsg, err, domain.ErrEmailNotVerified)
		}

		return nil, err
	}

	return &a, nil
}

func (r *accountRepository) UpdatePassword(ctx context.Context, accountID, newPassword string) (*domain.Account, error) {
	const errMsg = "accountRepository - UpdatePassword"

	sql, args, err := r.Builder.
		Update("account").
		Set("password", newPassword).
		Where(sq.Eq{"id": accountID}).
		Suffix(accountReturnAll).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	a, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[domain.Account])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.New(errMsg, err, domain.ErrAccountNotFound)
		}

		return nil, err
	}

	return &a, nil
}

func (r *accountRepository) SetPassword(ctx context.Context, token, newPassword string) (*domain.Account, error) {
	const errMsg = "accountRepository - SetPassword"

	sql := fmt.Sprintf(`
UPDATE account 
SET password = $1 
WHERE id = (
SELECT account_id 
FROM password_token
WHERE token = $2 
AND expires_at < now()
) %s`, accountReturnAll)

	rows, err := r.Pool.Query(ctx, sql, newPassword, token)
	if err != nil {
		return nil, err
	}

	a, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[domain.Account])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.New(errMsg, err, domain.ErrPasswordTokenExpired)
		}

		return nil, err
	}

	return &a, nil
}
