package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	"github.com/answersuck/vault/internal/domain/session"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

const sessionTable = "session"

type sessionPSQL struct {
	log    logging.Logger
	client *postgres.Client
}

func NewSessionPSQL(l logging.Logger, c *postgres.Client) *sessionPSQL {
	return &sessionPSQL{
		log:    l,
		client: c,
	}
}

func (r *sessionPSQL) Save(ctx context.Context, s *session.Session) error {
	sql := fmt.Sprintf(`
		INSERT INTO %s (id, account_id, max_age, user_agent, ip, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, sessionTable)

	r.log.Info("psql - session - Save: %s", sql)

	_, err := r.client.Pool.Exec(ctx, sql,
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
				return fmt.Errorf("psql - r.client.Pool.QueryRow.Scan: %w", session.ErrAlreadyExist)
			case pgerrcode.ForeignKeyViolation:
				return fmt.Errorf("psql - r.client.Pool.QueryRow.Scan: %w", session.ErrAccountNotFound)
			}
		}

		return fmt.Errorf("psql - r.client.Pool.QueryRow.Scan: %w", err)
	}

	return nil
}

func (r *sessionPSQL) FindById(ctx context.Context, sessionId string) (*session.Session, error) {
	sql := fmt.Sprintf(`
		SELECT
			s.account_id,
			s.user_agent,
			s.ip,
			s.max_age,
			s.expires_at,
			s.created_at
		FROM %s s
		WHERE s.id = $1
	`, sessionTable)

	r.log.Info("psql - session - FindById: %s", sql)

	s := session.Session{Id: sessionId}

	if err := r.client.Pool.QueryRow(ctx, sql, sessionId).Scan(
		&s.AccountId,
		&s.UserAgent,
		&s.IP,
		&s.MaxAge,
		&s.ExpiresAt,
		&s.CreatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("r.client.Pool.QueryRow.Scan: %w", session.ErrNotFound)
		}

		return nil, fmt.Errorf("r.client.Pool.QueryRow.Scan: %w", err)
	}

	return &s, nil
}

func (r *sessionPSQL) FindByIdWithVerified(ctx context.Context, sessionId string) (*session.SessionWithVerified, error) {
	sql := fmt.Sprintf(`
		SELECT
			s.account_id,
			s.user_agent,
			s.ip,
			s.max_age,
			s.expires_at,
			s.created_at,
			a.is_verified
		FROM %s s
		INNER JOIN %s a
		ON s.account_id = a.id
		WHERE s.id = $1
	`, sessionTable, accountTable)

	r.log.Info("psql - session - FindByIdWithVerified: %s", sql)

	s := session.SessionWithVerified{
		Session: session.Session{
			Id: sessionId,
		},
	}

	err := r.client.Pool.QueryRow(ctx, sql, sessionId).Scan(
		&s.Session.AccountId,
		&s.Session.UserAgent,
		&s.Session.IP,
		&s.Session.MaxAge,
		&s.Session.ExpiresAt,
		&s.Session.CreatedAt,
		&s.AccountVerified,
	)
	if err != nil {

		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("r.client.Pool.QueryRow.Scan: %w", session.ErrNotFound)
		}

		return nil, fmt.Errorf("r.client.Pool.QueryRow.Scan: %w", err)
	}

	return &s, nil
}

func (r *sessionPSQL) FindAll(ctx context.Context, accountId string) ([]*session.Session, error) {
	sql := fmt.Sprintf(`
		SELECT
			id,
			account_id,
			max_age,
			user_agent,
   			ip,
   			expires_at,
			created_at
		FROM %s
		WHERE account_id = $1
	`, sessionTable)

	r.log.Info("psql - session - FindAll: %s", sql)

	rows, err := r.client.Pool.Query(ctx, sql, accountId)
	if err != nil {
		return nil, fmt.Errorf("r.client.Pool.QueryRow.Scan: %w", err)
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
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		sessions = append(sessions, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return sessions, nil
}

func (r *sessionPSQL) Delete(ctx context.Context, sessionId string) error {
	sql := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, sessionTable)

	r.log.Info("psql - session - Delete: %s", sql)

	ct, err := r.client.Pool.Exec(ctx, sql, sessionId)
	if err != nil {
		return fmt.Errorf("r.client.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("r.client.Pool.Exec: %w", session.ErrNotDeleted)
	}

	return nil
}

func (r *sessionPSQL) DeleteWithExcept(ctx context.Context, accountId, sessionId string) error {
	sql := fmt.Sprintf(`DELETE FROM %s WHERE account_id = $1 AND id != $2`, sessionTable)

	r.log.Info("psql - session - DeleteWithExcept: %s", sql)

	ct, err := r.client.Pool.Exec(ctx, sql, accountId, sessionId)
	if err != nil {
		return fmt.Errorf("r.client.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("r.client.Pool.Exec: %w", session.ErrNotDeleted)
	}

	return nil
}

func (r *sessionPSQL) DeleteAll(ctx context.Context, accountId string) error {
	sql := fmt.Sprintf(`DELETE FROM %s WHERE account_id = $1`, sessionTable)

	r.log.Info("psql - session - DeleteAll: %s", sql)

	ct, err := r.client.Pool.Exec(ctx, sql, accountId)
	if err != nil {
		return fmt.Errorf("r.client.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("r.client.Pool.Exec: %w", session.ErrNotDeleted)
	}

	return nil
}
