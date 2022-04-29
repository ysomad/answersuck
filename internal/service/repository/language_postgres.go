package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/answersuck/vault/internal/domain"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

const (
	languageTable = "language"
)

type languageRepository struct {
	log logging.Logger
	*postgres.Client
}

func NewLanguageRepository(l logging.Logger, pg *postgres.Client) *languageRepository {
	return &languageRepository{l, pg}
}

func (r *languageRepository) FindAll(ctx context.Context) ([]*domain.Language, error) {
	sql := fmt.Sprintf(`
		SELECT id, name 
		FROM %s
	`, languageTable)

	r.log.Info("db: " + sql)

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("r.Pool.Query: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	defer rows.Close()

	var languages []*domain.Language

	for rows.Next() {
		var l domain.Language

		if err = rows.Scan(&l.Id, &l.Name); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", ErrNotFound)
		}

		languages = append(languages, &l)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return languages, nil
}
