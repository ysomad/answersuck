package pack

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (r *repository) GetWithTags(ctx context.Context, packID int32) (*entity.PackWithTags, error) {
	sql, args, err := r.Builder.
		Select(
			"p.id as id",
			"p.name as name",
			"p.author as author",
			"p.is_published as is_published",
			"p.cover_url as cover_url",
			"p.create_time as create_time",
			"pt.tag as tag",
		).
		From("packs p").
		LeftJoin("pack_tags pt on pt.pack_id = p.id").
		Where(squirrel.Eq{"id": packID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	pp, err := pgx.CollectRows(rows, pgx.RowToStructByName[packWithTag])
	if err != nil {
		return nil, fmt.Errorf("error getting pack with tags: %w", err)
	}

	if len(pp) == 0 {
		return nil, apperr.PackNotFound
	}

	pack := &entity.PackWithTags{
		Pack: entity.Pack{
			ID:         pp[0].ID,
			Name:       pp[0].Name,
			Author:     pp[0].Author,
			Published:  pp[0].Published,
			CoverURL:   string(pp[0].CoverURL),
			CreateTime: pp[0].CreateTime,
		},
		Tags: make([]string, len(pp)),
	}

	for i, p := range pp {
		pack.Tags[i] = p.Tag
	}

	return pack, nil
}
