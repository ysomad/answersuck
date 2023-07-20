package tag

import (
	"context"

	"github.com/ysomad/answersuck/internal/entity"
)

func (r *repository) Save(ctx context.Context, t entity.Tag) error {
	sql, args, err := r.Builder.
		Insert(tagTable).
		Columns("name", "author", "created_at").
		Values(t.Name, t.Author, t.CreatedAt).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}
