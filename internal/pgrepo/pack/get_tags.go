package pack

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (r *repository) GetTags(ctx context.Context, packID int32) ([]string, error) {
	sql, args, err := r.Builder.
		Select("tag").
		From(packTagsTable).
		Where(squirrel.Eq{"pack_id": packID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
		var tag string
		err := row.Scan(&tag)
		return tag, err
	})
}
