package media

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/ysomad/answersuck/internal/entity"
)

func (r *Repository) Save(ctx context.Context, m entity.Media) (entity.Media, error) {
	sql, args, err := r.Builder.
		Insert(mediaTable).
		Columns("url, type, uploader, create_time").
		Values(m.URL, m.Type, m.Uploader, m.CreateTime).
		Suffix("ON CONFLICT(url) DO UPDATE").
		Suffix("SET url = EXCLUDED.url").
		Suffix("RETURNING url, type, uploader, create_time").
		ToSql()
	if err != nil {
		return entity.Media{}, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return entity.Media{}, err
	}

	res, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[media])
	if err != nil {
		return entity.Media{}, err
	}

	return entity.Media(res), nil

}
