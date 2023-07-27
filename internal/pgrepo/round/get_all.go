package round

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/ysomad/answersuck/internal/entity"
)

func (r *repository) GetAll(ctx context.Context, packID int32) ([]entity.Round, error) {
	sql, args, err := r.Builder.
		Select("id, name, position").
		From(roundsTable).
		Where(squirrel.Eq{"pack_id": packID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Round])
}
