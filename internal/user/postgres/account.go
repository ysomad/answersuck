package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/ysomad/answersuck/internal/user/domain"
	"github.com/ysomad/answersuck/internal/user/service/dto"

	"github.com/ysomad/answersuck/pgclient"
)

type accountRepository struct {
	*pgclient.Client
}

func NewAccountRepository(c *pgclient.Client) *accountRepository {
	return &accountRepository{c}
}

func (r *accountRepository) Save(ctx context.Context, args dto.AccountSaveArgs) (*domain.Account, error) {
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
				return nil, domain.ErrAccountAlreadyExist
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

	return account.WithSaveArgs(args), nil
}

func (r *accountRepository) FindByID(ctx context.Context, accountID string) (*domain.Account, error) {
	return nil, nil
}

func (r *accountRepository) FindByEmail(ctx context.Context, email string) (*domain.Account, error) {
	return nil, nil
}

func (r *accountRepository) DeleteByID(ctx context.Context, accountID string) error {
	return nil
}
