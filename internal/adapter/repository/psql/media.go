package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/answersuck/host/internal/domain/media"
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
	sql := `
INSERT INTO media(id, filename, type, account_id, created_at)
VALUES ($1, $2, $3, $4, $5)`

	r.Debug("psql - media - Save", zap.String("sql", sql), zap.Any("media", m))

	_, err := r.Pool.Exec(
		ctx,
		sql,
		m.Id,
		m.Filename,
		m.Type,
		m.AccountId,
		m.CreatedAt,
	)
	if err != nil {
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

func (r *MediaRepo) FindMediaTypeById(ctx context.Context, mediaId string) (media.Type, error) {
	sql := "SELECT type FROM media WHERE id = $1"
	r.Debug("psql - media - FindMediaTypeById", zap.String("mediaId", mediaId))

	var mediaType media.Type
	err := r.Pool.QueryRow(ctx, sql, mediaId).Scan(&mediaType)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("psql - media - FindMediaTypeById - r.Pool.QueryRow.Scan: %w", media.ErrNotFound)
		}

		return "", fmt.Errorf("psql - media - FindMediaTypeById - r.Pool.QueryRow.Scan: %w", err)
	}

	return mediaType, nil
}
