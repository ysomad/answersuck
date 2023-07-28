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

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[round])
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, apperr.PackNotFound
	}

	rounds := make([]entity.Round, len(res))

	for i, round := range res {
		rounds[i] = entity.Round(round)
	}

	return rounds, nil
}
