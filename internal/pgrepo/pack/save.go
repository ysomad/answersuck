package pack

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype/zeronull"
	"github.com/ysomad/answersuck/internal/entity"
)

func (r *Repository) Save(ctx context.Context, p *entity.Pack) (int32, error) {
	sql, args, err := r.Builder.
		Insert(packTable).
		Columns("name, author, is_published, cover_url, create_time").
		Values(p.Name, p.Author, p.Published, zeronull.Text(p.CoverURL), p.CreateTime).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, err
	}

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&p.ID); err != nil {
		return 0, err
	}

	return p.ID, nil
}
