package roundtopic

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (r *Repository) Save(ctx context.Context, roundID, topicID int32) error {
	sql, args, err := r.Builder.
		Insert(roundTopicsTable).
		Columns("round_id, topic_id").
		Values(roundID, topicID).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.ConstraintName == "round_topics_pkey" {
			return apperr.TopicAlreadyInRound
		}

		return err
	}

	return nil
}
