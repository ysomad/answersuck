package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	"github.com/quizlyfun/quizly-backend/internal/domain"
	"github.com/quizlyfun/quizly-backend/pkg/postgres"
)

const accountTable = "account"

type accountRepository struct {
	*postgres.Postgres
}

func NewAccountRepository(pg *postgres.Postgres) *accountRepository {
	return &accountRepository{pg}
}

func (r *accountRepository) Create(ctx context.Context, acc domain.Account) (domain.Account, error) {
	sql, args, err := r.Builder.
		Insert(accountTable).
		Columns("username, email, password, is_verified").
		Values(acc.Username, acc.Email, acc.PasswordHash, acc.Verified).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return domain.Account{}, fmt.Errorf("r.Builder.Insert: %w", err)
	}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(&acc.ID); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.Account{}, fmt.Errorf("r.Pool.Exec: %w", domain.ErrAccountAlreadyExist)
			}
		}

		return domain.Account{}, fmt.Errorf("r.Pool.Exec: %w", err)
	}

	return acc, nil
}

func (r *accountRepository) FindByID(ctx context.Context, aid string) (domain.Account, error) {
	sql, args, err := r.Builder.
		Select("username, email, password, created_at, updated_at, is_verified").
		From(accountTable).
		Where(sq.Eq{"id": aid, "is_archived": false}).
		ToSql()
	if err != nil {
		return domain.Account{}, fmt.Errorf("r.Builder.Select: %w", err)
	}

	acc := domain.Account{ID: aid}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&acc.Username,
		&acc.Email,
		&acc.PasswordHash,
		&acc.CreatedAt,
		&acc.UpdatedAt,
		&acc.Verified,
	); err != nil {
		if err == pgx.ErrNoRows {
			return domain.Account{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", domain.ErrAccountNotFound)
		}

		return domain.Account{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return acc, nil
}

func (r *accountRepository) FindByEmail(ctx context.Context, email string) (domain.Account, error) {
	sql, args, err := r.Builder.
		Select("id, username, password, created_at, updated_at, is_verified").
		From(accountTable).
		Where(sq.Eq{"email": email, "is_archived": false}).
		ToSql()
	if err != nil {
		return domain.Account{}, fmt.Errorf("r.Builder.Select: %w", err)
	}

	a := domain.Account{Email: email}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&a.ID,
		&a.Username,
		&a.PasswordHash,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.Verified,
	); err != nil {
		if err == pgx.ErrNoRows {
			return domain.Account{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", domain.ErrAccountNotFound)
		}

		return domain.Account{}, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return a, nil
}

func (r *accountRepository) Archive(ctx context.Context, aid string, archive bool) error {
	sql, args, err := r.Builder.
		Update(accountTable).
		Set("is_archived", archive).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": aid, "is_archived": !archive}).
		ToSql()
	if err != nil {
		return fmt.Errorf("r.Builder.Update: %w", err)
	}

	ct, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("r.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("r.Pool.Exec: %w", domain.ErrAccountNotFound)
	}

	return nil
}
