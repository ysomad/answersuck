package psql

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/domain/tag"
	"go.uber.org/zap"

	"github.com/answersuck/vault/pkg/postgres"
)

type tagRepo struct {
	l *zap.Logger
	c *postgres.Client
}

func NewTagRepo(l *zap.Logger, c *postgres.Client) *tagRepo {
	return &tagRepo{
		l: l,
		c: c,
	}
}

func (r *tagRepo) SaveMultiple(ctx context.Context, req []tag.CreateReq) ([]*tag.Tag, error) {
	sql := `
		INSERT INTO tag(name, language_id)
		VALUES %s
		RETURNING id, name, language_id
	`

	argsNum := 2
	l := len(req)
	args := make([]any, 0, argsNum*l)

	for _, row := range req {
		args = append(args, row.Name, row.LanguageId)
	}

	sql = getBulkInsertSQLSimple(sql, argsNum, l)

	rows, err := r.c.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("psql - tag - SaveMultiple - r.c.Pool.Query: %w", err)
	}

	defer rows.Close()

	var tags []*tag.Tag

	for rows.Next() {
		var t tag.Tag

		if err = rows.Scan(&t.Id, &t.Name, &t.LanguageId); err != nil {
			return nil, fmt.Errorf("psql - tag - SaveMultiple - rows.Scan: %w", err)
		}

		tags = append(tags, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("psql - tag - SaveMultiple - rows.Err: %w", err)
	}

	return tags, nil
}

func (r *tagRepo) FindAll(ctx context.Context) ([]*tag.Tag, error) {
	sql := "SELECT id, name, language_id FROM tag"

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
