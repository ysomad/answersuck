package player

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/ysomad/answersuck/internal/entity"
)

func (r *Repository) GetOne(ctx context.Context, nickname string) (entity.Player, error) {
	sql, args, err := r.Builder.
		Select(
			"nickname",
			"email",
			"display_name",
			"email_verified",
			"password",
			"created_at",
			"updated_at",
		).
		From(playerTable).
		Where(squirrel.Eq{"nickname": nickname}).
		ToSql()
	if err != nil {
		return entity.Player{}, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return entity.Player{}, fmt.Errorf("r.Pool.Query: %w", err)
	}

	p, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[player])
	if err != nil {
		return entity.Player{}, fmt.Errorf("pgx.CollectOneRow: %w", err)
	}

	return entity.Player{
		Nickname:      p.Nickname,
		Email:         p.Email,
		DisplayName:   string(p.DisplayName),
		EmailVerified: p.EmailVerified,
		PasswordHash:  p.Password,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     time.Time(p.UpdatedAt),
	}, nil
}
