package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/answersuck/vault/internal/domain/session"

	"github.com/answersuck/vault/pkg/postgres"
)

type sessionRepo struct {
	l *zap.Logger
	c *postgres.Client
}

func NewSessionRepo(l *zap.Logger, c *postgres.Client) *sessionRepo {
	return &sessionRepo{
		l: l,
		c: c,
	}
}

func (r *sessionRepo) Save(ctx context.Context, s *session.Session) error {
	sql := `
		INSERT INTO session (id, account_id, max_age, user_agent, ip, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	r.l.Debug("psql - session - Save", zap.String("sql", sql), zap.Any("session", s))

	_, err := r.c.Pool.Exec(ctx, sql,
		s.Id,
		s.AccountId,
		s.MaxAge,
		s.UserAgent,
		s.IP,
		s.ExpiresAt,
		s.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return fmt.Errorf("psql - session - Save r.c.Pool.Exec: %w", session.ErrAlreadyExist)
			case pgerrcode.ForeignKeyViolation:
				return fmt.Errorf("psql - session - Save - r.c.Pool.Exec: %w", session.ErrAccountNotFound)
			}
		}

		return fmt.Errorf("psql - session - Save - r.c.Pool.Exec: %w", err)
	}

	return nil
}

func (r *sessionRepo) FindById(ctx context.Context, sessionId string) (*session.Session, error) {
	sql := `
		SELECT
			account_id,
			user_agent,
			ip,
			max_age,
			expires_at,
			created_at
		FROM session
		WHERE id = $1
	`

	r.l.Debug("psql - session - FindById", zap.String("sql", sql), zap.String("sessionId", sessionId))

	var s session.Session

	if err := r.c.Pool.QueryRow(ctx, sql, sessionId).Scan(
		&s.AccountId,
		&s.UserAgent,
		&s.IP,
		&s.MaxAge,
		&s.ExpiresAt,
		&s.CreatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("psql - session - FindById - r.c.Pool.QueryRow.Scan: %w", session.ErrNotFound)
		}

		return nil, fmt.Errorf("psql - session - FindById - r.c.Pool.QueryRow.Scan: %w", err)
	}

	s.Id = sessionId

	return &s, nil
}

func (r *sessionRepo) FindWithAccountDetails(ctx context.Context, sessionId string) (*session.WithAccountDetails, error) {
	sql := `
		SELECT
			s.account_id,
			s.user_agent,
			s.ip,
			s.max_age,
			s.expires_at,
			s.created_at,
			a.is_verified
		FROM session s
		INNER JOIN account a
		ON s.account_id = a.id
		WHERE s.id = $1
	`

	r.l.Debug("psql - session - FindWithAccountDetails", zap.String("sql", sql), zap.String("sessionId", sessionId))

	var (
		s session.Session
		d session.WithAccountDetails
	)

	err := r.c.Pool.QueryRow(ctx, sql, sessionId).Scan(
		&s.AccountId,
		&s.UserAgent,
		&s.IP,
		&s.MaxAge,
		&s.ExpiresAt,
		&s.CreatedAt,
		&d.Verified,
	)
	if err != nil {

		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("psql - session - FindWithAccountDetails - r.c.Pool.QueryRow.Scan: %w", session.ErrNotFound)
		}

		return nil, fmt.Errorf("psql - session - FindWithAccountDetails - r.c.Pool.QueryRow.Scan: %w", err)
	}

	d.Session = s

	return &d, nil
}

func (r *sessionRepo) FindAll(ctx context.Context, accountId string) ([]*session.Session, error) {
	sql := `
		SELECT
			id,
			account_id,
			max_age,
			user_agent,
   			ip,
   			expires_at,
			created_at
		FROM session
		WHERE account_id = $1
	`

	r.l.Debug("psql - session - FindAll", zap.String("sql", sql), zap.String("accountId", accountId))

	rows, err := r.c.Pool.Query(ctx, sql, accountId)
	if err != nil {
		return nil, fmt.Errorf("psql - session - FindAll - r.c.Pool.Query: %w", err)
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

func (r *sessionRepo) Delete(ctx context.Context, sessionId string) error {
	sql := "DELETE FROM session WHERE id = $1"

	r.l.Debug("psql - session - Delete", zap.String("sql", sql), zap.String("sessionId", sessionId))

	ct, err := r.c.Pool.Exec(ctx, sql, sessionId)
	if err != nil {
		return fmt.Errorf("psql - session - Delete - r.c.Pool.Exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("psql - session - Delete - r.c.Pool.Exec: %w", session.ErrNotDeleted)
	}

	return nil
}

func (r *sessionRepo) DeleteWithExcept(ctx context.Context, accountId, sessionId string) error {
	sql := "DELETE FROM session WHERE account_id = $1 AND id != $2"

	r.l.Debug(
		"psql - session - DeleteWithExcept",
		zap.String("sql", sql),
		zap.String("sessionId", sessionId),
		zap.String("accountId", accountId),
	)

	ct, err := r.c.Pool.Exec(ctx, sql, accountId, sessionId)
	if err != nil {
		return fmt.Errorf("psql - session - DeleteWithExcept - r.c.Pool.Exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("psql - session - DeleteWithExcept - r.c.Pool.Exec: %w", session.ErrNotDeleted)
	}

	return nil
}

func (r *sessionRepo) DeleteAll(ctx context.Context, accountId string) error {
	sql := "DELETE FROM session WHERE account_id = $1"
	r.l.Debug("psql - session - DeleteAll", zap.String("sql", sql), zap.String("accountId", accountId))

	ct, err := r.c.Pool.Exec(ctx, sql, accountId)
	if err != nil {
		return fmt.Errorf("psql - session - DeleteAll - r.c.Pool.Exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("psql - session - DeleteAll - r.c.Pool.Exec: %w", session.ErrNotDeleted)
	}

	return nil
}
