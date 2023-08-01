package roundtopic

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/postgres/topic"
)

func (r *Repository) GetAll(ctx context.Context, roundID int32) ([]entity.Topic, error) {
	sql, args, err := r.Builder.
		Select(
			"t.id as id",
			"t.title as title",
			"t.author as author",
			"t.create_time as create_time").
		From(roundTopicsTable + " rt").
		InnerJoin(topic.TopicsTable + " t ON rt.topic_id = t.id").
		Where(squirrel.Eq{"rt.round_id": roundID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	tt, err := pgx.CollectRows(rows, pgx.RowToStructByName[topic.Topic])
	if err != nil {
		return nil, err
	}

	topics := make([]entity.Topic, len(tt))

	for i, t := range tt {
		topics[i] = entity.Topic(t)
	}

	return topics, nil
}
