package tag

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (r *repository) Save(ctx context.Context, t entity.Tag) error {
	sql, args, err := r.Builder.
		Insert(TagsTable).
		Columns("name", "author", "create_time").
		Values(t.Name, t.Author, t.CreateTime).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && (pgErr.ConstraintName == "tags_pkey" || pgErr.ConstraintName == "players_email_key") {
			return apperr.TagAlreadyExists
		}

		return err
	}

	return nil
}
