package psql

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/answersuck/host/internal/domain/language"
	"github.com/answersuck/host/internal/pkg/postgres"
)

type languageRepo struct {
	*zap.Logger
	*postgres.Client
}

func NewLanguageRepo(l *zap.Logger, c *postgres.Client) *languageRepo {
	return &languageRepo{l, c}
}

func (r *languageRepo) FindAll(ctx context.Context) ([]language.Language, error) {
	sql := "SELECT id, name FROM language"
	r.Debug("psql - language - FindAll", zap.String("sql", sql))

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("psql - language - FindAll - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var langs []language.Language
	for rows.Next() {
		var l language.Language

		if err = rows.Scan(&l.Id, &l.Name); err != nil {
			return nil, fmt.Errorf("psql - language - FindAll - rows.Scan: %w", err)
		}

		langs = append(langs, l)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("psql - language - FindAll - rows.Err: %w", err)
	}

	return langs, nil
}
