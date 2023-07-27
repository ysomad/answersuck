package pack

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype/zeronull"
	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pgrepo/tag"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"golang.org/x/exp/slog"
)

func (r *repository) Save(ctx context.Context, p *entity.Pack, tags []string) (packID int32, err error) {
	if len(tags) == 0 {
		return r.insertPack(ctx, r.Pool, p)
	}

	err = pgx.BeginTxFunc(ctx, r.Pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		// 1. Insert tags
		if err = r.insertTags(ctx, tx, tags, p.Author, p.CreateTime); err != nil {
			return fmt.Errorf("error saving tags: %w", err)
		}

		// 2. Insert pack
		packID, err = r.insertPack(ctx, tx, p)
		if err != nil {
			return fmt.Errorf("error saving pack: %w", err)
		}

		// 3. Insert pack tags
		if err = r.insertPackTags(ctx, tx, packID, tags); err != nil {
			return fmt.Errorf("error saving pack tags: %w", err)
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return packID, nil
}

type queryRower interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func (r *repository) insertPack(ctx context.Context, db queryRower, p *entity.Pack) (int32, error) {
	sql, args, err := r.Builder.
		Insert(PacksTable).
		Columns("name, author, is_published, cover_url, create_time").
		Values(p.Name, p.Author, false, zeronull.Text(p.CoverURL), p.CreateTime).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, err
	}

	var packID int32

	if err := db.QueryRow(ctx, sql, args...).Scan(&packID); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.ConstraintName == "packs_cover_url_fkey" {
			return 0, apperr.MediaNotFound
		}

		return 0, err
	}

	return packID, nil
}

func (r *repository) insertTags(ctx context.Context, tx pgx.Tx, tags []string, author string, createTime time.Time) error {
	b := r.Builder.
		Insert(tag.TagsTable).
		Columns("name, author, create_time").
		Suffix("ON CONFLICT DO NOTHING")

	for _, t := range tags {
		b = b.Values(t, author, createTime)
	}

	sql, args, err := b.ToSql()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}

func (r *repository) insertPackTags(ctx context.Context, tx pgx.Tx, packID int32, tags []string) error {
	b := r.Builder.
		Insert(packTagsTable).
		Columns("pack_id, tag")

	for _, t := range tags {
		b = b.Values(packID, t)
	}

	sql, args, err := b.ToSql()
	if err != nil {
		return err
	}

	slog.Info(sql)
	slog.Info("pack", packID)

	if _, err := tx.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}
