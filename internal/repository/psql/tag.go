package repository

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/domain"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

const tagTable = "tag"

type tagRepo struct {
	log    logging.Logger
	client *postgres.Client
}

func NewTagRepo(l logging.Logger, c *postgres.Client) *tagRepo {
	return &tagRepo{
		log:    l,
		client: c,
	}
}

func (r *tagRepo) FindAll(ctx context.Context) ([]*domain.Tag, error) {
	sql := fmt.Sprintf(`SELECT id, name, language_id FROM %s`, tagTable)

	r.log.Info("psql - tag - FindAll: %s", sql)

	rows, err := r.client.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("r.client.Pool.QueryRow.Scan: %w", err)
	}

	defer rows.Close()

	var tags []*domain.Tag

	for rows.Next() {
		var t domain.Tag

		if err = rows.Scan(&t.Id, &t.Name, &t.LanguageId); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		tags = append(tags, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return tags, nil
}
