package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"

	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/pkg/postgres"
)

const (
	sessionTable = "session"
)

type sessionRepository struct {
	*postgres.Client
}

func NewSessionRepository(pg *postgres.Client) *sessionRepository {
	return &sessionRepository{pg}
}

func (r *sessionRepository) Create(ctx context.Context, s *domain.Session) (*domain.Session, error) {
	sql := fmt.Sprintf(`
		INSERT INTO %s (id, account_id, max_age, user_agent, ip, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, sessionTable)

	if err := r.Pool.QueryRow(ctx, sql,
		s.Id,
		s.AccountId,
		s.MaxAge,
		s.UserAgent,
		s.IP,
		s.ExpiresAt,
		s.CreatedAt,
	).Scan(&s.Id); err != nil {
		if err = isUniqueViolation(err); err != nil {
			return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
		}

		return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	return s, nil
}

func (r *sessionRepository) FindById(ctx context.Context, sid string) (*domain.Session, error) {
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

	s := domain.Session{Id: sid}

	if err := r.Pool.QueryRow(ctx, sql, sid).Scan(
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

func (r *sessionRepository) FindAll(ctx context.Context, aid string) ([]*domain.Session, error) {
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

	rows, err := r.Pool.Query(ctx, sql, aid)
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

func (r *sessionRepository) Delete(ctx context.Context, sid string) error {
	panic("implement")

	return nil
}

func (r *sessionRepository) DeleteAll(ctx context.Context, aid, sid string) error {
	panic("implement")

	return nil
}
