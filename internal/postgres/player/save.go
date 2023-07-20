package player

import (
	"context"

	"github.com/ysomad/answersuck/internal/entity"
)

func (r *Repository) Save(ctx context.Context, p entity.Player) error {
	sql, args, err := r.Builder.
		Insert(playerTable).
		Columns("nickname", "email", "password", "created_at").
		Values(p.Nickname, p.Email, p.PasswordHash, p.CreatedAt).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}
