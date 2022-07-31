package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/answersuck/vault/internal/domain/media"

	"github.com/answersuck/vault/pkg/postgres"
)

const (
	mediaTable = "media"
)

type mediaRepo struct {
	l *zap.Logger
	c *postgres.Client
}

func NewMediaRepo(l *zap.Logger, c *postgres.Client) *mediaRepo {
	return &mediaRepo{
		l: l,
		c: c,
	}
}

func (r *mediaRepo) Save(ctx context.Context, m media.Media) (media.Media, error) {
	sql := fmt.Sprintf(`
		INSERT INTO %s(id, url, type, account_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, mediaTable)

	_, err := r.c.Pool.Exec(
		ctx,
		sql,
		m.Id,
		m.URL,
		m.Type,
		m.AccountId,
		m.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return media.Media{}, fmt.Errorf("psql - media - Save - r.c.Pool.Exec: %w", media.ErrAlreadyExist)
			case pgerrcode.ForeignKeyViolation:
				return media.Media{}, fmt.Errorf("psql - media - Save - r.c.Pool.Exec: %w", media.ErrAccountNotFound)
			}
		}

		return media.Media{}, fmt.Errorf("psql - media - Save - r.c.Pool.Exec: %w", err)
	}

	return m, nil
}

func (r *mediaRepo) FindMimeTypeById(ctx context.Context, mediaId string) (string, error) {
	sql := fmt.Sprintf(`SELECT type FROM %s WHERE id = $1`, mediaTable)

	var mimeType string

	err := r.c.Pool.QueryRow(ctx, sql, mediaId).Scan(&mimeType)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("psql - media - FindMimeTypeById - r.c.Pool.QueryRow.Scan: %w", media.ErrNotFound)
		}

		return "", fmt.Errorf("psql - media - FindMimeTypeById - r.c.Pool.QueryRow.Scan: %w", err)
	}

	return mimeType, nil
}
