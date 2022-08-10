package psql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/answersuck/host/internal/domain/player"
	"github.com/answersuck/host/internal/pkg/postgres"
)

type PlayerRepo struct {
	*zap.Logger
	*postgres.Client
}

func NewPlayerRepo(l *zap.Logger, c *postgres.Client) *PlayerRepo {
	return &PlayerRepo{l, c}
}

func (r *PlayerRepo) FindByNickname(ctx context.Context, nickname string) (player.Player, error) {
	sql, args, err := r.Builder.
		Select("p.id, pa.filename").
		From("player p").
		InnerJoin("account a ON a.id = p.account_id").
		LeftJoin("player_avatar pa ON p.id = pa.player_id").
		Where(sq.And{
			sq.Eq{"a.nickname": nickname},
			sq.Eq{"a.is_archived": false},
		}).
		ToSql()
	if err != nil {
		return player.Player{}, fmt.Errorf("psql - player - FindByNickname - ToSql: %w", err)
	}

	r.Debug("psql - player - FindByNickname", zap.String("sql", sql), zap.Any("args", args))

	p := player.Player{Nickname: nickname}
	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&p.Id, &p.AvatarURL); err != nil {
		if err == pgx.ErrNoRows {
			return player.Player{}, fmt.Errorf("psql - player - FindByNickname - r.c.Pool.QueryRow.Scan: %w", player.ErrNotFound)
		}

		return player.Player{}, fmt.Errorf("psql - player - FindByNickname - r.c.Pool.QueryRow.Scan: %w", err)
	}

	return p, nil
}
