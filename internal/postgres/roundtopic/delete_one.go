package roundtopic

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (r *Repository) DeleteOne(ctx context.Context, roundID, topicID int32) error {
	sql, args, err := r.Builder.
		Delete(roundTopicsTable).
		Where(squirrel.And{
			squirrel.Eq{"round_id": roundID},
			squirrel.Eq{"topic_id": topicID},
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
		return apperr.RoundTopicNotDeleted
	}

	return nil
}
