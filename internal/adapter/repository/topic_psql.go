package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/answersuck/vault/internal/domain/topic"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

const topicTable = "topic"

type topicPSQL struct {
	log    logging.Logger
	client *postgres.Client
}

func NewTopicPSQL(l logging.Logger, c *postgres.Client) *topicPSQL {
	return &topicPSQL{
		log:    l,
		client: c,
	}
}

func (r *topicPSQL) Create(ctx context.Context, t topic.Topic) (int, error) {
	sql := fmt.Sprintf(`
		INSERT INTO %s(name, language_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, topicTable)

	r.log.Info("psql - topic - Create: %s", sql)

	var topicId int

	if err := r.client.Pool.QueryRow(ctx, sql,
		t.Name,
		t.LanguageId,
		t.CreatedAt,
		t.UpdatedAt,
	).Scan(&topicId); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return 0, fmt.Errorf("psql - r.client.Pool.QueryRow.Scan: %w", topic.ErrLanguageNotFound)
			}
		}

		return 0, fmt.Errorf("r.client.Pool.QueryRow.Scan: %w", err)
	}

	return topicId, nil
}

func (r *topicPSQL) FindAll(ctx context.Context) ([]*topic.Topic, error) {
	sql := fmt.Sprintf(`
		SELECT 
			id, 
			name, 
			language_id, 
			created_at, 
			updated_at 
		FROM %s
	`, topicTable)

	r.log.Info("psql - topic - FindAll: %s", sql)

	rows, err := r.client.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("r.client.Pool.QueryRow.Scan: %w", err)
	}

	defer rows.Close()

	var topics []*topic.Topic

	for rows.Next() {
		var t topic.Topic

		if err = rows.Scan(&t.Id, &t.Name, &t.LanguageId, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		topics = append(topics, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return topics, nil
}
