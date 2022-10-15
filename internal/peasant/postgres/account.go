package postgres

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/ysomad/answersuck/apperror"
	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"

	"github.com/ysomad/answersuck/pgclient"
)

type accountRepository struct {
	*pgclient.Client
}

func NewAccountRepository(c *pgclient.Client) *accountRepository {
	return &accountRepository{c}
}

func (r *accountRepository) Create(ctx context.Context, args dto.AccountSaveArgs) (*domain.Account, error) {
	const errMsg = "accountRepository - Create"

	accountSQL, accountArgs, err := r.Builder.
		Insert("account").
		Columns("email, username, password").
		Values(args.Email, args.Username, args.EncodedPassword).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var account domain.Account

	if err := tx.QueryRow(ctx, accountSQL, accountArgs...).Scan(
		&account.ID,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				switch pgErr.ConstraintName {
				case constraintAccountEmail:
					return nil, apperror.New(errMsg, err, domain.ErrEmailTaken)
				case constraintAccountUsername:
					return nil, apperror.New(errMsg, err, domain.ErrUsernameTaken)
				}

				return nil, err
			}
		}

		return nil, err
	}

	verifSQL, verifArgs, err := r.Builder.
		Insert("email_verification").
		Columns("code, account_id").
		Values(args.EmailVerifCode, account.ID).
		ToSql()
	if err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, verifSQL, verifArgs...); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	account.Email, account.Username = args.Email, args.Username
	return &account, nil
}

func (r *accountRepository) GetByID(ctx context.Context, accountID string) (*domain.Account, error) {
	const errMsg = "accountRepository - GetByID"

	sql, args, err := r.Builder.
		Select("id, username, email, is_email_verified, password, is_archived, created_at, updated_at").
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

	var account domain.Account

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&account.Username,
		&account.Email,
		&account.CreatedAt,
		&account.UpdatedAt,
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

	account.ID, account.Email = accountID, newEmail
	return &account, nil
}
