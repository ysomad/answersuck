package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"

	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

const topicTable = "topic"

type topic struct {
	log    logging.Logger
	client *postgres.Client
}

func NewTopic(l logging.Logger, c *postgres.Client) *topic {
	return &topic{
		log:    l,
		client: c,
	}
}

func (r *topic) Create(ctx context.Context, t domain.Topic) (int, error) {
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
		if err = unwrapError(err); err != nil {
			return 0, fmt.Errorf("r.client.Pool.QueryRow.Scan: %w", err)
		}

		return 0, fmt.Errorf("r.client.Pool.QueryRow.Scan: %w", err)
	}

	return topicId, nil
}

func (r *topic) FindAll(ctx context.Context) ([]*domain.Topic, error) {
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
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("r.client.Pool.Query: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("r.client.Pool.QueryRow.Scan: %w", err)
	}

	defer rows.Close()

	var topics []*domain.Topic

	for rows.Next() {
		var t domain.Topic

		if err = rows.Scan(&t.Id, &t.Name, &t.LanguageId, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", ErrNotFound)
		}

		topics = append(topics, &t)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return topics, nil
}
