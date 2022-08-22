package psql

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"

	"github.com/ysomad/answersuck-backend/internal/domain/tag"
	"github.com/ysomad/answersuck-backend/internal/pkg/filter"
	"github.com/ysomad/answersuck-backend/internal/pkg/pagination"
	"github.com/ysomad/answersuck-backend/internal/pkg/postgres"
	"github.com/ysomad/answersuck-backend/internal/pkg/sort"
)

type TagRepo struct {
	*zap.Logger
	*postgres.Client
}

func NewTagRepo(l *zap.Logger, c *postgres.Client) *TagRepo {
	return &TagRepo{l, c}
}

func (r *TagRepo) SaveMultiple(ctx context.Context, tags []tag.Tag) ([]tag.Tag, error) {
	if len(tags) < 1 {
		return nil, tag.ErrEmptyTagList
	}

	ib := r.Builder.
		Insert("tag").
		Columns("name, language_id")

	for _, t := range tags {
		ib = ib.Values(t.Name, t.LanguageId)
	}

	sql, args, err := ib.Suffix("RETURNING id, name, language_id").ToSql()
	if err != nil {
		return nil, fmt.Errorf("psql - tag - SaveMultiple - ToSql: %w", err)
	}

	r.Debug("psql - tag - SaveMultiple - ToSql", zap.String("sql", sql), zap.Any("args", args))

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return nil, fmt.Errorf("psql - tag - SaveMultiple - r.Pool.Query: %w", tag.ErrLanguageIdNotFound)
			}
		}

		return nil, fmt.Errorf("psql - tag - SaveMultiple - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var tagList []tag.Tag

	for rows.Next() {
		var t tag.Tag

		if err = rows.Scan(&t.Id, &t.Name, &t.LanguageId); err != nil {
			return nil, fmt.Errorf("psql - tag - SaveMultiple - rows.Scan: %w", err)
		}

		tagList = append(tagList, t)
	}

	if err = rows.Err(); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return nil, fmt.Errorf("psql - tag - SaveMultiple - rows.Err: %w", tag.ErrLanguageIdNotFound)
			}
		}

		return nil, fmt.Errorf("psql - tag - SaveMultiple - rows.Err: %w", err)
	}

	return tagList, nil
}

func (r *TagRepo) FindAll(ctx context.Context, p tag.ListParams) (pagination.List[tag.Tag], error) {
	sb := r.Builder.
		Select("id, name, language_id").
		From("tag").
		Where(sq.Gt{"id": p.Pagination.LastId})

	switch {
	case p.Filter.Name != "":
		sb = filter.New("name", filter.TypeILike, p.Filter.Name).UseSelectBuilder(sb)
	case p.Filter.LanguageId != 0:
		sb = filter.New("language_id", filter.TypeEQ, p.Filter.LanguageId).UseSelectBuilder(sb)
	}

	sb = sort.New("id", "ASC").UseSelectBuilder(sb)

	sql, args, err := sb.Limit(p.Pagination.Limit + 1).ToSql()
	if err != nil {
		return pagination.List[tag.Tag]{}, fmt.Errorf("psql - tag - FindAll - ToSql: %w", err)
	}

	r.Debug("psql - tag - FindAll - ToSql", zap.String("sql", sql), zap.Any("args", args))

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return pagination.List[tag.Tag]{}, fmt.Errorf("psql - tag - FindAll - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	tags := make([]tag.Tag, 0, p.Pagination.Limit+1)

	for rows.Next() {
		var t tag.Tag

		if err = rows.Scan(&t.Id, &t.Name, &t.LanguageId); err != nil {
			return pagination.List[tag.Tag]{}, fmt.Errorf("psql - tag - FindAll - rows.Scan: %w", err)
		}

		tags = append(tags, t)
	}

	if err = rows.Err(); err != nil {
		return pagination.List[tag.Tag]{}, fmt.Errorf("psql - tag - FindAll - rows.Err: %w", err)
	}

	return pagination.NewList(tags, p.Pagination.Limit), nil
}
