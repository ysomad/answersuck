package psql

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"

	"github.com/ysomad/answersuck-backend/internal/domain/answer"
	"github.com/ysomad/answersuck-backend/internal/pkg/filter"
	"github.com/ysomad/answersuck-backend/internal/pkg/pagination"
	"github.com/ysomad/answersuck-backend/internal/pkg/postgres"
	"github.com/ysomad/answersuck-backend/internal/pkg/sort"
)

type AnswerRepo struct {
	*zap.Logger
	*postgres.Client
}

func NewAnswerRepo(l *zap.Logger, c *postgres.Client) *AnswerRepo {
	return &AnswerRepo{l, c}
}

func (r *AnswerRepo) Save(ctx context.Context, a answer.Answer) (answer.Answer, error) {
	sql, args, err := r.Builder.
		Insert("answer").
		Columns("text, media_id, language_id").
		Values(a.Text, a.MediaId, a.LanguageId).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return answer.Answer{}, fmt.Errorf("psql - answer - Save - ToSql: %w", err)
	}

	r.Debug("psql - answer - Save - ToSql", zap.String("sql", sql), zap.Any("args", args))

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&a.Id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return answer.Answer{}, fmt.Errorf("psql - answer - Save - Scan: %w", answer.ErrLanguageNotFound)
			}
		}

		return answer.Answer{}, fmt.Errorf("psql - answer - Save - Scan: %w", err)
	}

	return a, nil
}

func (r *AnswerRepo) FindAll(ctx context.Context, p answer.ListParams) (pagination.List[answer.Answer], error) {
	sb := r.Builder.
		Select("id, text, media_id, language_id").
		From("answer").
		Where(sq.Gt{"id": p.Pagination.LastId})

	switch {
	case p.Filter.Text != "":
		sb = filter.New("text", filter.TypeILike, p.Filter.Text).UseSelectBuilder(sb)
	case p.Filter.LanguageId != 0:
		sb = filter.New("language_id", filter.TypeEQ, p.Filter.LanguageId).UseSelectBuilder(sb)
	}

	sb = sort.New("id", "ASC").UseSelectBuilder(sb)

	sql, args, err := sb.Limit(p.Pagination.Limit + 1).ToSql()
	if err != nil {
		return pagination.List[answer.Answer]{}, fmt.Errorf("psql - answer - FindAll - ToSql: %w", err)
	}

	r.Debug("psql - answer - FindAll - ToSql", zap.String("sql", sql), zap.Any("args", args))

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return pagination.List[answer.Answer]{}, fmt.Errorf("psql - answer - FindAll - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	answers := make([]answer.Answer, 0, p.Pagination.Limit+1)

	for rows.Next() {
		var a answer.Answer

		if err = rows.Scan(&a.Id, &a.Text, &a.MediaId, &a.LanguageId); err != nil {
			return pagination.List[answer.Answer]{}, fmt.Errorf("psql - answer - FindAll - Scan: %w", err)
		}

		answers = append(answers, a)
	}

	if err = rows.Err(); err != nil {
		return pagination.List[answer.Answer]{}, fmt.Errorf("psql - answer - FindAll - rows.Err: %w", err)
	}

	return pagination.NewList(answers, p.Pagination.Limit), nil
}
