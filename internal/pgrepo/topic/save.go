package topic

import (
	"context"

	"github.com/ysomad/answersuck/internal/entity"
)

func (r *Repository) Save(ctx context.Context, t entity.Topic) (int32, error) {
	sql, args, err := r.Builder.
		Insert(TopicsTable).
		Columns("title, author, create_time").
		Values(t.Title, t.Author, t.CreateTime).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, err
	}

	var topicID int32

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&topicID); err != nil {
		return 0, err
	}

	return topicID, nil
}
