package psql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/answersuck/vault/internal/domain/player"
	"github.com/answersuck/vault/pkg/postgres"
)

type playerRepo struct {
	l *zap.Logger
	c *postgres.Client
}

func NewPlayerRepo(l *zap.Logger, c *postgres.Client) *playerRepo {
	return &playerRepo{
		l: l,
		c: c,
	}
}

func (r *playerRepo) FindByNickname(ctx context.Context, nickname string) (player.Player, error) {
	sql := `
		SELECT
			p.id,
			a.nickname,
			pa.url
		FROM player p
		INNER JOIN account a ON a.id = p.account_id
		LEFT JOIN player_avatar pa ON p.id = pa.player_id
		WHERE
			a.nickname = $1
			AND a.is_archived = $2
	`

	var p player.Player

	if err := r.c.Pool.QueryRow(ctx, sql, nickname, false).Scan(
		&p.Id,
		&p.Nickname,
		&p.AvatarURL,
	); err != nil {

		if err == pgx.ErrNoRows {
			return player.Player{}, fmt.Errorf("psql - player - FindByNickname - r.c.Pool.QueryRow.Scan: %w", player.ErrNotFound)
		}

		return player.Player{}, fmt.Errorf("psql - player - FindByNickname - r.c.Pool.QueryRow.Scan: %w", err)
	}

	return p, nil
}
