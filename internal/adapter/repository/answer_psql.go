package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/answersuck/vault/internal/domain/answer"
	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

const (
	answerTable = "answer"
)

type answerPSQL struct {
	log    logging.Logger
	client *postgres.Client
}

func NewAnswerPSQL(l logging.Logger, c *postgres.Client) *answerPSQL {
	return &answerPSQL{
		log:    l,
		client: c,
	}
}

func (r *answerPSQL) Save(ctx context.Context, a answer.Answer) (answer.Answer, error) {
	sql := fmt.Sprintf(`
		INSERT INTO %s(text, image)	
		VALUES ($1, $2)
		RETURNING id
	`, answerTable)

	r.log.Info("psql - answer: %s", sql)

	if err := r.client.Pool.QueryRow(ctx, sql, a.Text, a.MediaId).Scan(&a.Id); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return answer.Answer{}, fmt.Errorf("psql - r.client.Pool.QueryRow.Scan: %w", answer.ErrMediaNotFound)
			}

		}

		return answer.Answer{}, fmt.Errorf("r.client.Pool.QueryRow.Scan: %w", err)
	}

	return a, nil
}
