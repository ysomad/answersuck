package player

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (r *Repository) Save(ctx context.Context, p *entity.Player) error {
	sql, args, err := r.Builder.
		Insert(playerTable).
		Columns("nickname", "email", "password", "created_at").
		Values(p.Nickname, p.Email, p.PasswordHash, p.CreatedAt).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && (pgErr.ConstraintName == "players_pkey" || pgErr.ConstraintName == "players_email_key") {
			return apperr.ErrPlayerAlreadyExists
		}

		return err
	}

	return nil
}
