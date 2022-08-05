package psql

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/answersuck/host/internal/domain/tag"
	"github.com/answersuck/host/internal/pkg/postgres"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

type TagRepo struct {
	*zap.Logger
	*postgres.Client
}

func NewTagRepo(l *zap.Logger, c *postgres.Client) *TagRepo {
	return &TagRepo{l, c}
}

func (r *TagRepo) SaveMultiple(ctx context.Context, tags []tag.Tag) ([]tag.Tag, error) {
	if len(tags) <= 0 {
		return nil, tag.ErrEmptyTagList
	}

	insertBuilder := r.Builder.
		Insert("tag").
		Columns("name, language_id")

	for _, t := range tags {
		insertBuilder = insertBuilder.Values(t.Name, t.LanguageId)
	}

	sql, args, err := insertBuilder.Suffix("RETURNING id, name, language_id").ToSql()
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

func (r *TagRepo) FindAll(ctx context.Context) ([]tag.Tag, error) {
	sql, _, err := r.Builder.
		Select("id, name, language_id").
		From("tag").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("psql - tag - FindAll - ToSql: %w", err)
	}

	r.Debug("psql - tag - FindAll - ToSql", zap.String("sql", sql))

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("psql - tag - FindAll - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var tags []tag.Tag

	for rows.Next() {
		var t tag.Tag

		if err = rows.Scan(&t.Id, &t.Name, &t.LanguageId); err != nil {
			return nil, fmt.Errorf("psql - tag - FindAll - rows.Scan: %w", err)
		}

		tags = append(tags, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("psql - tag - FindAll - rows.Err: %w", err)
	}

	return tags, nil
}
