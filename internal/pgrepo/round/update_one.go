package round

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (r *Repository) UpdateOne(ctx context.Context, round entity.Round) error {
	sql, args, err := r.Builder.
		Update(RoundsTable).
		SetMap(map[string]interface{}{
			"name":     round.Name,
			"position": round.Position,
		}).
		Where(sq.And{
			sq.Eq{"id": round.ID},
			sq.Eq{"pack_id": round.PackID},
		}).
		ToSql()
	if err != nil {
		return err
	}

	ct, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return apperr.RoundNotFound
	}

	return nil
}
