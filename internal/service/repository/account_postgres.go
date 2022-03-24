package repository

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"github.com/quizlyfun/quizly-backend/internal/domain"
	"github.com/quizlyfun/quizly-backend/internal/dto"
	"github.com/quizlyfun/quizly-backend/pkg/postgres"
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

func (r *accountRepository) insertAccount(ctx context.Context, tx pgx.Tx, a *domain.Account) (*domain.Account, error) {
	sql, args, err := r.Builder.
		Insert(accountTable).
		Columns("username, email, password, is_verified").
		Values(a.Username, a.Email, a.PasswordHash, a.Verified).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("r.Builder.Insert: %w", err)
	}

	if err = tx.QueryRow(ctx, sql, args...).Scan(&a.Id); err != nil {

		if err := isUniqueViolation(err); err != nil {
			return nil, fmt.Errorf("r.Pool.Exec: %w", err)
		}

		return nil, fmt.Errorf("r.Pool.Exec: %w", err)
	}

	return a, nil
}

func (r *accountRepository) insertAccountVerification(ctx context.Context, tx pgx.Tx, code, aid string) error {
	sql, args, err := r.Builder.
		Insert(accountVerificationTable).
		Columns("code, account_id").
		Values(code, aid).
		ToSql()

	if err != nil {
		return fmt.Errorf("r.Builder.Insert: %w", err)
	}

	if _, err = tx.Exec(ctx, sql, args...); err != nil {

		if err := isUniqueViolation(err); err != nil {
			return fmt.Errorf("r.Pool.Exec: %w", err)
		}

		return fmt.Errorf("r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *accountRepository) insertAccountAvatar(ctx context.Context, tx pgx.Tx, url, aid string) error {
	sql, args, err := r.Builder.
		Insert(accountAvatarTable).
		Columns("url, account_id").
		Values(url, aid).
		ToSql()

	if err != nil {
		return fmt.Errorf("r.Builder.Insert: %w", err)
	}

	if _, err = tx.Exec(ctx, sql, args...); err != nil {

		if err := isUniqueViolation(err); err != nil {
			return fmt.Errorf("r.Pool.Exec: %w", err)
		}

		return fmt.Errorf("r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *accountRepository) Create(ctx context.Context, a *domain.Account) (*domain.Account, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("r.Pool.Begin: %w", err)
	}

	defer tx.Rollback(ctx)

	a, err = r.insertAccount(ctx, tx, a)
	if err != nil {
		return nil, fmt.Errorf("r.insertAccount: %w", err)
	}

	if err = r.insertAccountVerification(ctx, tx, a.VerificationCode, a.Id); err != nil {
		return nil, fmt.Errorf("r.insertAccountVerification: %w", err)
	}

	if err = r.insertAccountAvatar(ctx, tx, a.AvatarURL, a.Id); err != nil {
		return nil, fmt.Errorf("r.insertAccountAvatar: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("tx.Commit: %w", err)
	}

	return a, nil
}

func (r *accountRepository) FindById(ctx context.Context, aid string) (*domain.Account, error) {
	sql, args, err := r.Builder.
		Select(
			"a.username",
			"a.email",
			"a.password",
			"a.created_at",
			"a.updated_at",
			"a.is_verified",
			"av.url as avatar_url",
		).
		From(fmt.Sprintf("%s a", accountTable)).
		LeftJoin(fmt.Sprintf("%s av ON av.account_id = a.id", accountAvatarTable)).
		Where(sq.Eq{"a.id": aid, "a.is_archived": false}).
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
	sql, args, err := r.Builder.
		Select(
			"a.id",
			"a.username",
			"a.password",
			"a.created_at",
			"a.updated_at",
			"a.is_verified",
			"av.url as avatar_url",
		).
		From(fmt.Sprintf("%s a", accountTable)).
		LeftJoin(fmt.Sprintf("%s av ON av.account_id = a.id", accountAvatarTable)).
		Where(sq.Eq{"a.email": email, "a.is_archived": false}).
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
	sql, args, err := r.Builder.
		Select(
			"a.id",
			"a.email",
			"a.password",
			"a.created_at",
			"a.updated_at",
			"a.is_verified",
			"av.url as avatar_url",
		).
		From(fmt.Sprintf("%s a", accountTable)).
		LeftJoin(fmt.Sprintf("%s av ON av.account_id = a.id", accountAvatarTable)).
		Where(sq.Eq{"a.username": username, "a.is_archived": false}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("r.Builder.Select: %w", err)
	}

	a := domain.Account{Username: username}

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
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
		return fmt.Errorf("r.Pool.Exec: %w", ErrNoAffectedRows)
	}

	return nil
}

func (r *accountRepository) Verify(ctx context.Context, a dto.AccountVerification) error {
	sql := `
		UPDATE account a
		SET 
			is_verified = $1,
			updated_at = $2
		FROM account_verification av
		WHERE
			a.is_verified = $3
			AND av.account_id = $4
			AND av.code = $5
	`

	ct, err := r.Pool.Exec(ctx, sql, a.Verified, a.UpdatedAt, !a.Verified, a.AccountId, a.Code)
	if err != nil {
		return fmt.Errorf("r.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("r.Pool.Exec: %w", ErrNoAffectedRows)
	}

	return nil
}
