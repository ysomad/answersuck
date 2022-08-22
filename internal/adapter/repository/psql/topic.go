package psql

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"

	"github.com/ysomad/answersuck-backend/internal/domain/topic"
	"github.com/ysomad/answersuck-backend/internal/pkg/filter"
	"github.com/ysomad/answersuck-backend/internal/pkg/pagination"
	"github.com/ysomad/answersuck-backend/internal/pkg/postgres"
	"github.com/ysomad/answersuck-backend/internal/pkg/sort"
)

type TopicRepo struct {
	*zap.Logger
	*postgres.Client
}

func NewTopicRepo(l *zap.Logger, c *postgres.Client) *TopicRepo {
	return &TopicRepo{l, c}
}

func (r *TopicRepo) Save(ctx context.Context, t topic.Topic) (topic.Topic, error) {
	sql, args, err := r.Builder.
		Insert("topic").
		Columns("name, language_id, created_at").
		Values(t.Name, t.LanguageId, t.CreatedAt).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return topic.Topic{}, fmt.Errorf("psql - topic - Save - ToSql: %w", err)
	}

	r.Debug("psql - topic - Save", zap.String("sql", sql), zap.Any("args", args))

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&t.Id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return topic.Topic{}, fmt.Errorf("psql - topic - Save - Scan: %w", topic.ErrLanguageNotFound)
			}
		}

		return topic.Topic{}, fmt.Errorf("psql - topic - Save - Scan: %w", err)
	}

	return t, nil
}

func (r *TopicRepo) FindAll(ctx context.Context, p topic.ListParams) (pagination.List[topic.Topic], error) {
	sb := r.Builder.
		Select("id, name, language_id, created_at").
		From("topic").
		Where(sq.Gt{"id": p.Pagination.LastId})

	switch {
	case p.Filter.Name != "":
		sb = filter.New("name", filter.TypeILike, p.Filter.Name).UseSelectBuilder(sb)
	case p.Filter.LanguageId != 0:
		sb = filter.New("language_id", filter.TypeEQ, p.Filter.LanguageId).UseSelectBuilder(sb)
	}

	sb = sort.New("id", "ASC").UseSelectBuilder(sb)

	sql, args, err := sb.Limit(p.Pagination.Limit + 1).ToSql()
	if err != nil {
		return pagination.List[topic.Topic]{}, fmt.Errorf("psql - topic - FindAll - ToSql: %w", err)
	}

	r.Debug("psql - topic - Save", zap.String("sql", sql), zap.Any("args", args))

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return pagination.List[topic.Topic]{}, fmt.Errorf("psql - topic - FindAll - Scan: %w", err)
	}
	defer rows.Close()

	var topics []topic.Topic
	for rows.Next() {
		var t topic.Topic

		if err = rows.Scan(&t.Id, &t.Name, &t.LanguageId, &t.CreatedAt); err != nil {
			return pagination.List[topic.Topic]{}, fmt.Errorf("psql - topic - FindAll - rows.Scan: %w", err)
		}

		topics = append(topics, t)
	}

	if err = rows.Err(); err != nil {
		return pagination.List[topic.Topic]{}, fmt.Errorf("psql - topic - FindAll - rows.Err: %w", err)
	}

	return pagination.NewList(topics, p.Pagination.Limit), nil
}
