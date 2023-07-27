package session

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ Store = &Postgres{}

type Postgres struct {
	*pgxpool.Pool
}

func NewPostgresStore(p *pgxpool.Pool) *Postgres {
	return &Postgres{p}
}

func (p *Postgres) Save(ctx context.Context, s *Session) error {
	if _, err := p.Exec(
		ctx,
		"INSERT INTO sessions (id, user_agent, player_ip, player_verified, player_nickname, expire_time) VALUES ($1, $2, $3, $4, $5, $6)",
		s.ID, s.User.UserAgent, s.User.IP, s.User.Verified, s.User.ID, s.ExpiresAt,
	); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) Get(ctx context.Context, sid string) (*Session, error) {
	var s Session

	err := p.
		QueryRow(ctx, "SELECT id, user_agent, player_ip, player_verified, player_nickname, expire_time FROM sessions WHERE id = $1", sid).
		Scan(&s.ID, &s.User.UserAgent, &s.User.IP, &s.User.Verified, &s.User.ID, &s.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return &s, nil
}

func (p *Postgres) Delete(ctx context.Context, sid string) error {
	_, err := p.Exec(ctx, "DELETE FROM sessions WHERE id = ?", sid)

	return err
}
