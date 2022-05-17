package repository

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/domain/language"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

const languageTable = "language"

type languagePSQL struct {
	log    logging.Logger
	client *postgres.Client
}

func NewLanguagePSQL(l logging.Logger, c *postgres.Client) *languagePSQL {
	return &languagePSQL{
		log:    l,
		client: c,
	}
}

func (r *languagePSQL) FindAll(ctx context.Context) ([]*language.Language, error) {
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

	var langs []*language.Language

	for rows.Next() {
		var l language.Language

		if err = rows.Scan(&l.Id, &l.Name); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		langs = append(langs, &l)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return langs, nil
}
