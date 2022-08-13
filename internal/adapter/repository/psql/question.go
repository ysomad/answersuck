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

	"github.com/answersuck/host/internal/domain/question"
	"github.com/answersuck/host/internal/pkg/filter"
	"github.com/answersuck/host/internal/pkg/pagination"
	"github.com/answersuck/host/internal/pkg/postgres"
	"github.com/answersuck/host/internal/pkg/sort"
)

type QuestionRepo struct {
	*zap.Logger
	*postgres.Client
}

func NewQuestionRepo(l *zap.Logger, c *postgres.Client) *QuestionRepo {
	return &QuestionRepo{l, c}
}

func (r *QuestionRepo) Save(ctx context.Context, dto question.CreateDTO) (uint32, error) {
	sql, args, err := r.Builder.
		Insert("question").
		Columns("text, answer_id, account_id, media_id, language_id, created_at").
		Values(dto.Text, dto.AnswerId, dto.AccountId, dto.MediaId, dto.LanguageId, dto.CreatedAt).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("psql - question - Save - ToSql: %w", err)
	}

	r.Debug("psql - question - Save", zap.String("sql", sql), zap.Any("args", args))

	var questionId uint32

	if err = r.Pool.QueryRow(ctx, sql, args...).Scan(&questionId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return 0, fmt.Errorf("psql - question - Save - Scan: %w", question.ErrForeignKeyViolation)
			}
		}

		return 0, fmt.Errorf("psql - question - Save - Scan: %w", err)
	}

	return questionId, nil
}

func (r *QuestionRepo) FindById(ctx context.Context, questionId uint32) (question.Detailed, error) {
	sql, args, err := r.Builder.
		Select("q.text",
			"answer.text",
			"answer_media.filename",
			"answer_media.type",
			"account.nickname",
			"question_media.filename",
			"question_media.type",
			"q.language_id",
			"q.created_at").
		From("question q").
		InnerJoin("account on account.id = q.account_id").
		InnerJoin("answer on answer.id = q.answer_id").
		LeftJoin("media question_media on question_media.id = q.media_id").
		LeftJoin("media answer_media on answer_media.id = answer.media_id").
		Where(sq.Eq{"q.id": questionId}).
		ToSql()
	if err != nil {
		return question.Detailed{}, fmt.Errorf("psql - question - FindById - ToSql: %w", err)
	}

	r.Debug("psql - question - FindById", zap.String("sql", sql), zap.Any("args", args))

	q := question.Detailed{Id: questionId}
	if err = r.Pool.QueryRow(ctx, sql, questionId).Scan(
		&q.Text,
		&q.Answer,
		&q.AnswerMediaURL,
		&q.AnswerMediaType,
		&q.Author,
		&q.MediaURL,
		&q.MediaType,
		&q.LanguageId,
		&q.CreatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return question.Detailed{}, fmt.Errorf("psql - question - FindById - Scan: %w", question.ErrNotFound)
		}

		return question.Detailed{}, fmt.Errorf("psql - question - FindById - Scan: %w", err)
	}

	return q, nil
}

func (r *QuestionRepo) FindAll(ctx context.Context, p question.ListParams) (pagination.List[question.Minimized], error) {
	var sb sq.SelectBuilder

	sb = r.Builder.
		Select("q.id, q.text, q.language_id, q.created_at, a.nickname").
		From("question q").
		InnerJoin("account a ON a.id = q.account_id").
		Where(sq.Gt{"q.id": p.Pagination.LastId})

	switch {
	case p.Filter.Text != "":
		sb = filter.New("q.text", filter.TypeILike, p.Filter.Text).UseSelectBuilder(sb)
	case p.Filter.Author != "":
		sb = filter.New("a.nickname", filter.TypeEQ, p.Filter.Author).UseSelectBuilder(sb)
	case p.Filter.LanguageId != 0:
		sb = filter.New("q.language_id", filter.TypeEQ, p.Filter.LanguageId).UseSelectBuilder(sb)
	}

	sb = sort.New("id", "ASC").UseSelectBuilder(sb)

	sql, args, err := sb.Limit(p.Pagination.Limit + 1).ToSql()
	if err != nil {
		return pagination.List[question.Minimized]{}, fmt.Errorf("psql - question - FindAll - ToSql: %w", err)
	}

	r.Debug("psql - question - FindAll - ToSql", zap.String("sql", sql), zap.Any("args", args))

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return pagination.List[question.Minimized]{}, fmt.Errorf("psql - question - FindAll - Query: %w", err)
	}
	defer rows.Close()

	var qs []question.Minimized
	for rows.Next() {
		var q question.Minimized

		if err = rows.Scan(&q.Id, &q.Text, &q.LanguageId, &q.CreatedAt, &q.Author); err != nil {
			return pagination.List[question.Minimized]{}, fmt.Errorf("psql - question - FindAll - rows.Scan: %w", err)
		}

		qs = append(qs, q)
	}

	if rows.Err() != nil {
		return pagination.List[question.Minimized]{}, fmt.Errorf("psql - question - FindAll - rows.Err: %w", err)
	}

	return pagination.NewList(qs, p.Pagination.Limit), nil
}
