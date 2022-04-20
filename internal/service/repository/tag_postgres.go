package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"

	"github.com/answersuck/vault/internal/domain"

	"github.com/answersuck/vault/pkg/postgres"
)

const (
	tagTable = "tag"
)

type tagRepository struct {
	*postgres.Client
}

func NewTagRepository(pg *postgres.Client) *tagRepository {
	return &tagRepository{pg}
}

func (r *tagRepository) FindAll(ctx context.Context) ([]*domain.Tag, error) {
	sql := fmt.Sprintf(`SELECT id, name, language_id FROM %s`, tagTable)

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("r.Pool.Query: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	defer rows.Close()

	var tags []*domain.Tag

	for rows.Next() {
		var t domain.Tag

		if err = rows.Scan(&t.Id, &t.Name, &t.LanguageId); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", ErrNotFound)
		}

		tags = append(tags, &t)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return tags, nil
}
