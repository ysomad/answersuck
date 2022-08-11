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
		Select("p.id, p.avatar_filename, a.is_verified").
		From("player p").
		InnerJoin("account a ON a.id = p.account_id").
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
	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&p.Id, &p.AvatarFilename, &p.Verified); err != nil {
		if err == pgx.ErrNoRows {
			return player.Player{}, fmt.Errorf("psql - player - FindByNickname - r.c.Pool.QueryRow.Scan: %w", player.ErrNotFound)
		}

		return player.Player{}, fmt.Errorf("psql - player - FindByNickname - r.c.Pool.QueryRow.Scan: %w", err)
	}

	return p, nil
}

func (r *PlayerRepo) SetAvatar(ctx context.Context, a player.Avatar) error {
	sql, args, err := r.Builder.
		Update("player").
		Set("avatar_filename", a.Filename).
		Set("updated_at", a.UpdatedAt).
		Where(sq.Eq{"account_id": a.AccountId}).
		ToSql()
	if err != nil {
		return fmt.Errorf("psql - player - SetAvatar - ToSql: %w", err)
	}

	r.Debug("psql - player - SetAvatar", zap.String("sql", sql), zap.Any("args", args))

	ct, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("psql - player - SetAvatar - Exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("psql - player - SetAvatar - Exec: %w", player.ErrNotFound)
	}

	return nil
}
