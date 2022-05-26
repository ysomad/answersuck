package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	"github.com/answersuck/vault/internal/domain/media"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

const (
	mediaTable = "media"
)

type mediaPSQL struct {
	log    logging.Logger
	client *postgres.Client
}

func NewMediaPSQL(l logging.Logger, c *postgres.Client) *mediaPSQL {
	return &mediaPSQL{
		log:    l,
		client: c,
	}
}

func (r *mediaPSQL) Save(ctx context.Context, m media.Media) (media.Media, error) {
	sql := fmt.Sprintf(`
		INSERT INTO %s(id, url, mime_type, account_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, mediaTable)

	r.log.Info("psql - media - Save: %s", sql)

	_, err := r.client.Pool.Exec(
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
				return media.Media{}, fmt.Errorf("psql - r.client.Pool.Query: %w", media.ErrAlreadyExist)
			case pgerrcode.ForeignKeyViolation:
				return media.Media{}, fmt.Errorf("psql - r.client.Pool.Query: %w", media.ErrAccountNotFound)
			}
		}

		return media.Media{}, fmt.Errorf("psql - r.client.Pool.Query: %w", err)
	}

	return m, nil
}

func (r *mediaPSQL) FindMimeTypeById(ctx context.Context, mediaId string) (string, error) {
	sql := fmt.Sprintf(`SELECT mime_type FROM %s WHERE id = $1`, mediaTable)

	r.log.Info("psql - media - FindMimeTypeById: %s", sql)

	var mimeType string

	err := r.client.Pool.QueryRow(ctx, sql, mediaId).Scan(&mimeType)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("psql - r.client.Pool.QueryRow.Scan: %w", media.ErrNotFound)
		}

		return "", fmt.Errorf("psql - r.client.Pool.QueryRow.Scan: %w", err)
	}

	return mimeType, nil
}
