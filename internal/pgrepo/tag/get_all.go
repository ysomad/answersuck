package tag

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/paging"
	"github.com/ysomad/answersuck/internal/pkg/sort"
)

func (r *repository) GetAll(ctx context.Context, p paging.Params, sorts []sort.Sort) (paging.List[entity.Tag], error) {
	limit, offset, err := paging.OffsetToken(p.PageToken).Decode()
	if err != nil {
		return paging.List[entity.Tag]{}, err
	}

	// use limit from params only if token has no limit
	if limit == 0 {
		limit = uint64(p.PageSize)
	}

	b := r.Builder.
		Select("name, author, create_time").
		From(TagsTable).
		Limit(limit + 1).
		Offset(offset)

	for _, sort := range sorts {
		b = sort.Attach(b)
	}

	sql, args, err := b.ToSql()
	if err != nil {
		return paging.List[entity.Tag]{}, fmt.Errorf("b.ToSql: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return paging.List[entity.Tag]{}, fmt.Errorf("r.Pool.Query: %w", err)
	}

	tags, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Tag])
	if err != nil {
		return paging.List[entity.Tag]{}, fmt.Errorf("pgx.RowToStructPyBos: %w", err)
	}

	return paging.NewListWithOffset(tags, limit, offset)
}
