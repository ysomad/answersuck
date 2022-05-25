package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/answersuck/vault/internal/domain/media"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

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

	_, err := r.client.Pool.Query(
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
