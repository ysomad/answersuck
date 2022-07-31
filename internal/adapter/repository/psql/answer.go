package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"

	"github.com/answersuck/vault/internal/domain/answer"
	"github.com/answersuck/vault/pkg/postgres"
)

const (
	answerTable = "answer"
)

type answerRepo struct {
	l *zap.Logger
	c *postgres.Client
}

func NewAnswerRepo(l *zap.Logger, c *postgres.Client) *answerRepo {
	return &answerRepo{
		l: l,
		c: c,
	}
}

func (r *answerRepo) Save(ctx context.Context, a answer.Answer) (answer.Answer, error) {
	sql := fmt.Sprintf(`
		INSERT INTO %s(text, image)
		VALUES ($1, $2)
		RETURNING id
	`, answerTable)

	if err := r.c.Pool.QueryRow(ctx, sql, a.Text, a.MediaId).Scan(&a.Id); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return answer.Answer{}, fmt.Errorf("psql - answer - Save - r.c.Pool.QueryRow.Scan: %w", answer.ErrMediaNotFound)
			}

		}

		return answer.Answer{}, fmt.Errorf("psql - answer - Save - r.c.Pool.QueryRow.Scan: %w", err)
	}

	return a, nil
}
