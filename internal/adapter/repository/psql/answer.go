package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"

	"github.com/answersuck/host/internal/domain/answer"
	"github.com/answersuck/host/internal/pkg/postgres"
)

type answerRepo struct {
	*zap.Logger
	*postgres.Client
}

func NewAnswerRepo(l *zap.Logger, c *postgres.Client) *answerRepo {
	return &answerRepo{l, c}
}

func (r *answerRepo) Save(ctx context.Context, a answer.Answer) (answer.Answer, error) {
	sql, args, err := r.Builder.
		Insert("answer").
		Columns("text, media_id").
		Values(a.Text, a.MediaId).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return answer.Answer{}, fmt.Errorf("psql - answer - Save - ToSql: %w", err)
	}

	r.Debug("psql - answer - Save - ToSql", zap.String("sql", sql), zap.Any("args", args))

	if err := r.Pool.QueryRow(ctx, sql, a.Text, a.MediaId).Scan(&a.Id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return answer.Answer{}, fmt.Errorf("psql - answer - Save - Scan: %w", answer.ErrMediaNotFound)
			}
		}

		return answer.Answer{}, fmt.Errorf("psql - answer - Save - Scan: %w", err)
	}

	return a, nil
}
