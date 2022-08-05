package psql

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/answersuck/host/internal/domain/session"
	"github.com/answersuck/host/internal/pkg/postgres"
)

type SessionRepo struct {
	*zap.Logger
	*postgres.Client
}

func NewSessionRepo(l *zap.Logger, c *postgres.Client) *SessionRepo {
	return &SessionRepo{l, c}
}

func (r *SessionRepo) Save(ctx context.Context, s *session.Session) error {
	sql, args, err := r.Builder.
		Insert("session").
		Columns("id, account_id, max_age, user_agent, ip, expires_at, created_at").
		Values(s.Id, s.AccountId, s.MaxAge, s.UserAgent, s.IP, s.ExpiresAt, s.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("psql - session - Save - ToSql: %w", err)
	}

	r.Debug("psql - session - Save", zap.String("sql", sql), zap.Any("args", args))

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return fmt.Errorf("psql - session - Save - r.Pool.Exec: %w", session.ErrAlreadyExist)
			case pgerrcode.ForeignKeyViolation:
				return fmt.Errorf("psql - session - Save - r.Pool.Exec: %w", session.ErrAccountNotFound)
			}
		}

		return fmt.Errorf("psql - session - Save - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *SessionRepo) FindById(ctx context.Context, sessionId string) (*session.Session, error) {
	sql, args, err := r.Builder.
		Select("account_id, user_agent, ip, max_age, expires_at, created_at").
		From("session").
		Where(sq.Eq{"id": sessionId}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("psql - session - FindById - ToSql: %w", err)
	}

	r.Debug("psql - session - FindById", zap.String("sql", sql), zap.Any("args", args))

	s := session.Session{Id: sessionId}

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&s.AccountId,
		&s.UserAgent,
		&s.IP,
		&s.MaxAge,
		&s.ExpiresAt,
		&s.CreatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("psql - session - FindById - r.Pool.QueryRow.Scan: %w", session.ErrNotFound)
		}

		return nil, fmt.Errorf("psql - session - FindById - r.Pool.QueryRow.Scan: %w", err)
	}

	s.Id = sessionId

	return &s, nil
}

func (r *SessionRepo) FindWithAccountDetails(ctx context.Context, sessionId string) (*session.WithAccountDetails, error) {
	sql, args, err := r.Builder.
		Select("s.account_id, s.user_agent, s.ip, s.max_age, s.expires_at, s.created_at, a.is_verified").
		From("session s").
		InnerJoin("account a ON s.account_id = a.id").
		Where(sq.Eq{"s.id": sessionId}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("psql - session - FindWithAccountDetails - ToSql: %w", err)
	}

	r.Debug("psql - session - FindWithAccountDetails", zap.String("sql", sql), zap.Any("args", args))

	var (
		s session.Session
		d session.WithAccountDetails
	)
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&s.AccountId,
		&s.UserAgent,
		&s.IP,
		&s.MaxAge,
		&s.ExpiresAt,
		&s.CreatedAt,
		&d.AccountVerified,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("psql - session - FindWithAccountDetails - r.Pool.QueryRow.Scan: %w", session.ErrNotFound)
		}

		return nil, fmt.Errorf("psql - session - FindWithAccountDetails - r.Pool.QueryRow.Scan: %w", err)
	}

	s.Id = sessionId
	d.Session = s

	return &d, nil
}

func (r *SessionRepo) FindAll(ctx context.Context, accountId string) ([]*session.Session, error) {
	sql, args, err := r.Builder.
		Select("id, account_id, max_age, user_agent, ip, expires_at, created_at").
		From("session").
		Where(sq.Eq{"account_id": accountId}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("psql - session - FindAll - ToSql: %w", err)
	}

	r.Debug("psql - session - FindAll", zap.String("sql", sql), zap.Any("args", args))

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("psql - session - FindAll - r.Pool.Query: %w", err)
	}

	defer rows.Close()

	var sessions []*session.Session

	for rows.Next() {
		var s session.Session

		if err = rows.Scan(
			&s.Id,
			&s.AccountId,
			&s.MaxAge,
			&s.UserAgent,
			&s.IP,
			&s.ExpiresAt,
			&s.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("psql - session - FindAll - rows.Scan: %w", err)
		}

		sessions = append(sessions, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("psql - session - FindAll - rows.Err: %w", err)
	}

	return sessions, nil
}

func (r *SessionRepo) Delete(ctx context.Context, sessionId string) error {
	sql, args, err := r.Builder.Delete("session").Where(sq.Eq{"id": sessionId}).ToSql()
	if err != nil {
		return fmt.Errorf("psql - session - Delete - ToSql: %w", err)
	}

	r.Debug("psql - session - Delete", zap.String("sql", sql), zap.Any("args", args))

	ct, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("psql - session - Delete - r.Pool.Exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("psql - session - Delete - r.Pool.Exec: %w", session.ErrNotFound)
	}

	return nil
}

func (r *SessionRepo) DeleteAllWithExcept(ctx context.Context, accountId, sessionId string) error {
	sql := "DELETE FROM session WHERE account_id = $1 AND id != $2"
	sql, args, err := r.Builder.
		Delete("session").
		Where(sq.And{
			sq.Eq{"account_id": accountId},
			sq.NotEq{"id": sessionId},
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("psql - session - DeleteAllWithExcept - ToSql: %w", err)
	}

	r.Debug("psql - session - DeleteAllWithExcept", zap.Any("args", args))

	ct, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("psql - session - DeleteAllWithExcept - r.Pool.Exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("psql - session - DeleteAllWithExcept - r.Pool.Exec: %w", session.ErrAccountNotFound)
	}

	return nil
}

func (r *SessionRepo) DeleteAll(ctx context.Context, accountId string) error {
	sql, args, err := r.Builder.Delete("session").Where(sq.Eq{"account_id": accountId}).ToSql()
	if err != nil {
		return fmt.Errorf("psql - session - DeleteAll - ToSql: %w", err)
	}

	r.Debug("psql - session - DeleteAll", zap.String("sql", sql), zap.Any("args", args))

	ct, err := r.Pool.Exec(ctx, sql, accountId)
	if err != nil {
		return fmt.Errorf("psql - session - DeleteAll - r.Pool.Exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("psql - session - DeleteAll - r.Pool.Exec: %w", session.ErrAccountNotFound)
	}

	return nil
}
