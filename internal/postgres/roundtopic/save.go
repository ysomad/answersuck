package roundtopic

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (r *Repository) Save(ctx context.Context, roundID, topicID int32) (int32, error) {
	sql, args, err := r.Builder.
		Insert(roundTopicsTable).
		Columns("round_id, topic_id").
		Values(roundID, topicID).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, err
	}

	var id int32

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.ConstraintName == "round_topics_pkey" {
			return 0, apperr.RoundTopicAlreadyExists
		}

		return 0, err
	}

	return id, nil
}
