package round

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (r *Repository) GetAll(ctx context.Context, packID int32) ([]entity.Round, error) {
	sql, args, err := r.Builder.
		Select("id, name, pack_id, position").
		From(RoundsTable).
		Where(squirrel.Eq{"pack_id": packID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	rr, err := pgx.CollectRows(rows, pgx.RowToStructByName[round])
	if err != nil {
		return nil, err
	}

	if len(rr) == 0 {
		return nil, apperr.PackNotFound
	}

	rounds := make([]entity.Round, len(rr))

	for i, r := range rr {
		rounds[i] = entity.Round(r)
	}

	return rounds, nil
}
