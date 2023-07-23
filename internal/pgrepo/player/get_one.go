package player

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (r *Repository) GetOne(ctx context.Context, login string, loginType entity.LoginType) (*entity.Player, error) {
	b := r.Builder.
		Select(
			"nickname",
			"email",
			"display_name",
			"email_verified",
			"password",
			"create_time",
			"update_time",
		).
		From(playerTable)

	switch loginType {
	case entity.LoginTypeEmail:
		b = b.Where(sq.Eq{"email": login})
	default:
		b = b.Where(sq.Eq{"nickname": login})
	}

	sql, args, err := b.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("r.Pool.Query: %w", err)
	}

	p, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[player])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrPlayerNotFound
		}

		return nil, fmt.Errorf("pgx.CollectOneRow: %w", err)
	}

	return &entity.Player{
		Nickname:      p.Nickname,
		Email:         p.Email,
		DisplayName:   string(p.DisplayName),
		EmailVerified: p.EmailVerified,
		PasswordHash:  p.Password,
		CreatedAt:     p.CreateTime,
		UpdateTime:    time.Time(p.UpdateTime),
	}, nil
}
