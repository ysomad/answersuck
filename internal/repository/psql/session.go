package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

const sessionTable = "session"

type session struct {
	log    logging.Logger
	client *postgres.Client
}

func NewSession(l logging.Logger, c *postgres.Client) *session {
	return &session{
		log:    l,
		client: c,
	}
}

func (r *session) Create(ctx context.Context, s *domain.Session) (*domain.Session, error) {
	sql := fmt.Sprintf(`
		INSERT INTO %s (id, account_id, max_age, user_agent, ip, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, sessionTable)

	r.log.Info("psql - session - Create: %s", sql)

	if err := r.client.Pool.QueryRow(ctx, sql,
		s.Id,
		s.AccountId,
		s.MaxAge,
		s.UserAgent,
		s.IP,
		s.ExpiresAt,
		s.CreatedAt,
	).Scan(&s.Id); err != nil {
		if err = unwrapError(err); err != nil {
			return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
		}

		return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return s, nil
}

func (r *session) FindById(ctx context.Context, sid string) (*domain.Session, error) {
	sql := fmt.Sprintf(`
		SELECT 
			account_id,
			user_agent,
			ip,
			max_age,
			expires_at,
			created_at
		FROM %s
		WHERE id = $1
	`, sessionTable)

	r.log.Info("psql - session - FindById: %s", sql)

	s := domain.Session{Id: sid}

	if err := r.client.Pool.QueryRow(ctx, sql, sid).Scan(
		&s.AccountId,
		&s.UserAgent,
		&s.IP,
		&s.MaxAge,
		&s.ExpiresAt,
		&s.CreatedAt,
	); err != nil {

		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return &s, nil
}

func (r *session) FindAll(ctx context.Context, aid string) ([]*domain.Session, error) {
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

	rows, err := r.client.Pool.Query(ctx, sql, aid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("r.Pool.Query: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	defer rows.Close()

	var sessions []*domain.Session

	for rows.Next() {
		var s domain.Session

		if err = rows.Scan(
			&s.Id,
			&s.AccountId,
			&s.MaxAge,
			&s.UserAgent,
			&s.IP,
			&s.ExpiresAt,
			&s.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", ErrNotFound)
		}

		sessions = append(sessions, &s)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return sessions, nil
}

func (r *session) Delete(ctx context.Context, sid string) error {
	sql := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, sessionTable)

	r.log.Info("psql - session - Delete: %s", sql)

	ct, err := r.client.Pool.Exec(ctx, sql, sid)
	if err != nil {
		return fmt.Errorf("r.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("r.Pool.Exec: %w", ErrNoAffectedRows)
	}

	return nil
}

func (r *session) DeleteWithExcept(ctx context.Context, aid, sid string) error {
	sql := fmt.Sprintf(`DELETE FROM %s WHERE account_id = $1 AND id != $2`, sessionTable)

	r.log.Info("psql - session - DeleteWithExcept: %s", sql)

	ct, err := r.client.Pool.Exec(ctx, sql, aid, sid)
	if err != nil {
		return fmt.Errorf("r.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("r.Pool.Exec: %w", ErrNoAffectedRows)
	}

	return nil
}

func (r *session) DeleteAll(ctx context.Context, aid string) error {
	sql := fmt.Sprintf(`DELETE FROM %s WHERE account_id = $1`, sessionTable)

	r.log.Info("psql - session - DeleteAll: %s", sql)

	ct, err := r.client.Pool.Exec(ctx, sql, aid)
	if err != nil {
		return fmt.Errorf("r.Pool.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("r.Pool.Exec: %w", ErrNoAffectedRows)
	}

	return nil
}
