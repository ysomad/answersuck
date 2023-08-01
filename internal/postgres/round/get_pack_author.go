package round

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"github.com/ysomad/answersuck/internal/postgres/pack"
)

func (r *Repository) GetPackAuthor(ctx context.Context, roundID int32) (string, error) {
	sql, args, err := r.Builder.
		Select("p.author").
		From(RoundsTable + " r").
		InnerJoin(pack.PacksTable + " p ON r.pack_id = p.id").
		Where(squirrel.Eq{"r.id": roundID}).
		ToSql()
	if err != nil {
		return "", err
	}

	var author string

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&author); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", apperr.RoundNotFound
		}

		return "", err
	}

	return author, nil
}
