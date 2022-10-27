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

	table        string
	returningAll string
}

func NewAccountRepository(c *pgclient.Client) *accountRepository {
	return &accountRepository{
		c,
		"account",
		"RETURNING id, username, email, is_email_verified, is_archived, created_at, updated_at",
	}
}

func (r *accountRepository) Create(ctx context.Context, args dto.AccountCreateArgs) (*domain.Account, error) {
	const errMsg = "accountRepository - Create"

	query, queryArgs, err := r.Builder.
		Insert(r.table).
		Columns("email, username, password").
		Values(args.Email, args.Username, args.Password).
		Suffix(r.returningAll).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, query, queryArgs...)
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

	return &account, nil
}

func (r *accountRepository) GetByID(ctx context.Context, accountID string) (*domain.Account, error) {
	const errMsg = "accountRepository - GetByID"

	sql, args, err := r.Builder.
		Select("id, username, email, is_email_verified, is_archived, created_at, updated_at").
		From(r.table).
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
		Update(r.table).
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
		From(r.table).
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

	query, queryArgs, err := r.Builder.
		Update(r.table).
		Set("email", newEmail).
		Set("is_email_verified", false).
		Where(sq.And{
			sq.Eq{"id": accountID},
			sq.Eq{"is_archived": false},
		}).
		Suffix(r.returningAll).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	a, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[domain.Account])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.New(errMsg, err, domain.ErrAccountNotFound)
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == constraintAccountEmail {
			return nil, apperror.New(errMsg, pgErr, domain.ErrEmailTaken)
		}

		return nil, err
	}

	return &a, nil
}

func (r *accountRepository) VerifyEmail(ctx context.Context, accountID string) (*domain.Account, error) {
	const errMsg = "accountRepository - VerifyEmail"

	query, queryArgs, err := r.Builder.
		Update(r.table).
		Set("is_email_verified", true).
		Where(sq.And{
			sq.Eq{"id": accountID},
			sq.Eq{"is_email_verified": false},
		}).
		Suffix(r.returningAll).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	a, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[domain.Account])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.New(errMsg, err, domain.ErrEmailAlreadyVerified)
		}

		return nil, err
	}

	return &a, nil
}

func (r *accountRepository) UpdatePassword(ctx context.Context, accountID, newPassword string) (*domain.Account, error) {
	const errMsg = "accountRepository - UpdatePassword"

	query, queryArgs, err := r.Builder.
		Update(r.table).
		Set("password", newPassword).
		Where(sq.Eq{"id": accountID}).
		Suffix(r.returningAll).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, query, queryArgs...)
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
