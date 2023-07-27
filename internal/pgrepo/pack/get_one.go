package pack

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (r *repository) GetOne(ctx context.Context, packID int32) (*entity.Pack, error) {
	sql, args, err := r.Builder.
		Select("id, name, author, is_published, cover_url, create_time").
		From(PacksTable).
		Where(squirrel.Eq{"id": packID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	p, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[pack])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.PackNotFound
		}

		return nil, err
	}

	return &entity.Pack{
		ID:         p.ID,
		Name:       p.Name,
		Author:     p.Author,
		Published:  p.Published,
		CoverURL:   string(p.CoverURL),
		CreateTime: p.CreateTime,
	}, nil
}
