package psql

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/domain/tag"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

const tagTable = "tag"

type tagRepo struct {
	l logging.Logger
	c *postgres.Client
}

func NewTagRepo(l logging.Logger, c *postgres.Client) *tagRepo {
	return &tagRepo{
		l: l,
		c: c,
	}
}

func (r *tagRepo) FindAll(ctx context.Context) ([]*tag.Tag, error) {
	sql := fmt.Sprintf(`SELECT id, name, language_id FROM %s`, tagTable)

	r.l.Info("psql - tag - FindAll: %s", sql)

	rows, err := r.c.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("psql - tag - FindAll - r.c.Pool.Query: %w", err)
	}

	defer rows.Close()

	var tags []*tag.Tag

	for rows.Next() {
		var t tag.Tag

		if err = rows.Scan(&t.Id, &t.Name, &t.LanguageId); err != nil {
			return nil, fmt.Errorf("psql - tag - FindAll - rows.Scan: %w", err)
		}

		tags = append(tags, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("psql - tag - FindAll - rows.Err: %w", err)
	}

	return tags, nil
}
