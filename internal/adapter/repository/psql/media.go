package psql

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/answersuck/host/internal/domain/media"
	"github.com/answersuck/host/internal/pkg/mime"
	"github.com/answersuck/host/internal/pkg/postgres"
)

type MediaRepo struct {
	*zap.Logger
	*postgres.Client
}

func NewMediaRepo(l *zap.Logger, c *postgres.Client) *MediaRepo {
	return &MediaRepo{l, c}
}

func (r *MediaRepo) Save(ctx context.Context, m media.Media) error {
	sql, args, err := r.Builder.
		Insert("media").
		Columns("id, filename, type, account_id, created_at").
		Values(m.Id, m.Filename, m.Type, m.AccountId, m.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("psql - media - Save - ToSql: %w", err)
	}

	r.Debug("psql - media - Save", zap.String("sql", sql), zap.Any("args", args))

	if _, err = r.Pool.Exec(ctx, sql, args...); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return fmt.Errorf("psql - media - Save - r.Pool.Exec: %w", media.ErrAlreadyExist)
			case pgerrcode.ForeignKeyViolation:
				return fmt.Errorf("psql - media - Save - r.Pool.Exec: %w", media.ErrAccountNotFound)
			}
		}

		return fmt.Errorf("psql - media - Save - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *MediaRepo) FindById(ctx context.Context, mediaId string) (media.Media, error) {
	sql, args, err := r.Builder.
		Select("filename, type, account_id, created_at").
		From("media").
		Where(sq.Eq{"id": mediaId}).
		ToSql()
	if err != nil {
		return media.Media{}, fmt.Errorf("psql - media - FindById - ToSql: %w", err)
	}

	r.Debug("psql - media - FindById", zap.String("sql", sql), zap.Any("args", args))

	m := media.Media{Id: mediaId}
	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&m.Filename,
		&m.Type,
		&m.AccountId,
		&m.CreatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return media.Media{}, fmt.Errorf("psql - media - FindById - Scan: %w", media.ErrNotFound)
		}

		return media.Media{}, fmt.Errorf("psql - media - FindById - Scan: %w", err)
	}

	return m, nil
}

func (r *MediaRepo) FindMediaTypeById(ctx context.Context, mediaId string) (mime.Type, error) {
	sql, args, err := r.Builder.
		Select("type").
		From("media").
		Where(sq.Eq{"id": mediaId}).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("psql - media - FindMediaTypeById - ToSql: %w", err)
	}

	r.Debug("psql - media - FindMediaTypeById", zap.Any("args", args))

	var mediaType string
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&mediaType)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("psql - media - FindMediaTypeById - r.Pool.QueryRow.Scan: %w", media.ErrNotFound)
		}

		return "", fmt.Errorf("psql - media - FindMediaTypeById - r.Pool.QueryRow.Scan: %w", err)
	}

	return mime.Type(mediaType), nil
}
