package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"

	"github.com/answersuck/vault/internal/domain/topic"

	"github.com/answersuck/vault/pkg/postgres"
)

const topicTable = "topic"

type topicRepo struct {
	l *zap.Logger
	c *postgres.Client
}

func NewTopicRepo(l *zap.Logger, c *postgres.Client) *topicRepo {
	return &topicRepo{
		l: l,
		c: c,
	}
}

func (r *topicRepo) Save(ctx context.Context, t topic.Topic) (topic.Topic, error) {
	sql := fmt.Sprintf(`
		INSERT INTO %s(name, language_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, topicTable)

	if err := r.c.Pool.QueryRow(ctx, sql,
		t.Name,
		t.LanguageId,
		t.CreatedAt,
		t.UpdatedAt,
	).Scan(&t.Id); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return topic.Topic{}, fmt.Errorf("psql - topic - Save - r.c.Pool.QueryRow.Scan: %w", topic.ErrLanguageNotFound)
			}
		}

		return topic.Topic{}, fmt.Errorf("psql - topic - Save - r.c.Pool.QueryRow.Scan: %w", err)
	}

	return t, nil
}

func (r *topicRepo) FindAll(ctx context.Context) ([]*topic.Topic, error) {
	sql := fmt.Sprintf(`
		SELECT
			id,
			name,
			language_id,
			created_at,
			updated_at
		FROM %s
	`, topicTable)

	rows, err := r.c.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("psql - topic - FindAll - r.c.Pool.QueryRow.Scan: %w", err)
	}

	defer rows.Close()

	var topics []*topic.Topic

	for rows.Next() {
		var t topic.Topic

		if err = rows.Scan(&t.Id, &t.Name, &t.LanguageId, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("psql - topic - FindAll - rows.Scan: %w", err)
		}

		topics = append(topics, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("psql - topic - FindAll - rows.Err: %w", err)
	}

	return topics, nil
}
