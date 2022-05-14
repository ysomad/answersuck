package repository

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/domain"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

const languageTable = "language"

type languageRepo struct {
	log    logging.Logger
	client *postgres.Client
}

func NewLanguageRepo(l logging.Logger, c *postgres.Client) *languageRepo {
	return &languageRepo{
		log:    l,
		client: c,
	}
}

func (r *languageRepo) FindAll(ctx context.Context) ([]*domain.Language, error) {
	sql := fmt.Sprintf(`
		SELECT id, name 
		FROM %s
	`, languageTable)

	r.log.Info("psql - language - FindAll: %s", sql)

	rows, err := r.client.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("r.client.Pool.QueryRow.Scan: %w", err)
	}

	defer rows.Close()

	var languages []*domain.Language

	for rows.Next() {
		var l domain.Language

		if err = rows.Scan(&l.Id, &l.Name); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		languages = append(languages, &l)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return languages, nil
}
