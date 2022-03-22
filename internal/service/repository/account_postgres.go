package repository

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"github.com/quizlyfun/quizly-backend/internal/domain"
	"github.com/quizlyfun/quizly-backend/pkg/postgres"
)

const (
	accountTable = "account"
)

type accountRepository struct {
	*postgres.Postgres
}

func NewAccountRepository(pg *postgres.Postgres) *accountRepository {
	return &accountRepository{pg}
}

func (r *accountRepository) Create(ctx context.Context, a *domain.Account) (*domain.Account, error) {
	sql, args, err := r.Builder.
		Insert(accountTable).
		Columns("id, username, email, password, is_verified").
		Values(a.Id, a.Username, a.Email, a.PasswordHash, a.Verified).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("r.Builder.Insert: %w", err)
	}

	if _, err = r.Pool.Exec(ctx, sql, args...); err != nil {

		if err := isUniqueViolation(err); err != nil {
			return nil, fmt.Errorf("r.Pool.Exec: %w", err)
		}

		return nil, fmt.Errorf("r.Pool.Exec: %w", err)
	}

	return a, nil
}

func (r *accountRepository) FindByID(ctx context.Context, aid string) (*domain.Account, error) {
	sql, args, err := r.Builder.
		Select("username, email, password, created_at, updated_at, is_verified").
		From(accountTable).
		Where(sq.Eq{"id": aid, "is_archived": false}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("r.Builder.Select: %w", err)
	}

	a := domain.Account{Id: aid}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&a.Username,
		&a.Email,
		&a.PasswordHash,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.Verified,
	); err != nil {

		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return &a, nil
}

func (r *accountRepository) FindByEmail(ctx context.Context, email string) (*domain.Account, error) {
	sql, args, err := r.Builder.
		Select("id, username, password, created_at, updated_at, is_verified").
		From(accountTable).
		Where(sq.Eq{"email": email, "is_archived": false}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("r.Builder.Select: %w", err)
	}

	a := domain.Account{Email: email}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&a.Id,
		&a.Username,
		&a.PasswordHash,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.Verified,
	); err != nil {

		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return &a, nil
}

func (r *accountRepository) FindByUsername(ctx context.Context, u string) (*domain.Account, error) {
	sql, args, err := r.Builder.
		Select("id, email, password, created_at, updated_at, is_verified").
		From(accountTable).
		Where(sq.Eq{"username": u, "is_archived": false}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("r.Builder.Select: %w", err)
	}

	a := domain.Account{Username: u}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&a.Id,
		&a.Email,
		&a.PasswordHash,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.Verified,
	); err != nil {

		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return &a, nil
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
		return fmt.Errorf("r.Pool.Exec: %w", ErrNotFound)
	}

	return nil
}
