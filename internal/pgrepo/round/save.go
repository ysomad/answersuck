package round

import (
	"context"

	"github.com/ysomad/answersuck/internal/entity"
)

func (r *repository) Save(ctx context.Context, round entity.Round) (int32, error) {
	sql, args, err := r.Builder.
		Insert(roundsTable).
		Columns("name, position, pack_id").
		Values(round.Name, round.Position, round.PackID).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, err
	}

	var roundID int32

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(&roundID); err != nil {
		return 0, err
	}

	return roundID, nil
}
