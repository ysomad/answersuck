package repository

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/domain/tag"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

const tagTable = "tag"

type tagPSQL struct {
	log    logging.Logger
	client *postgres.Client
}

func NewTagPSQL(l logging.Logger, c *postgres.Client) *tagPSQL {
	return &tagPSQL{
		log:    l,
		client: c,
	}
}

func (r *tagPSQL) FindAll(ctx context.Context) ([]*tag.Tag, error) {
	sql := fmt.Sprintf(`SELECT id, name, language_id FROM %s`, tagTable)

	r.log.Info("psql - tag - FindAll: %s", sql)

	rows, err := r.client.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("r.client.Pool.QueryRow.Scan: %w", err)
	}

	defer rows.Close()

	var tags []*tag.Tag

	for rows.Next() {
		var t tag.Tag

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
