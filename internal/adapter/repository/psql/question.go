package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/answersuck/host/internal/domain/question"
	"github.com/answersuck/host/internal/pkg/postgres"
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

func (r *QuestionRepo) FindAll(ctx context.Context) ([]question.Minimized, error) {
	sql := "SELECT id, text, language_id FROM question"

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("psql - question - FindAll - r.c.Pool.Query: %w", err)
	}

	defer rows.Close()

	var qs []question.Minimized

	for rows.Next() {
		var q question.Minimized

		if err = rows.Scan(&q.Id, &q.Text, &q.LanguageId); err != nil {
			return nil, fmt.Errorf("psql - question - FindAll - rows.Scan: %w", err)
		}

		qs = append(qs, q)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("psql - question - FindAll - rows.Err: %w", err)
	}

	return qs, nil
}

func (r *QuestionRepo) FindById(ctx context.Context, questionId int) (*question.Detailed, error) {
	sql := `
		SELECT
			q.text,
			ans.text AS answer,
			am.url AS answer_image_url,
			acc.nickname AS author,
			qm.url AS media_url,
			qm.type AS media_type,
			q.language_id,
			q.created_at,
			q.updated_at
		FROM question q
		INNER JOIN account acc on acc.id = q.account_id
		INNER JOIN answer ans on ans.id = q.answer_id
		LEFT JOIN media qm on qm.id = q.media_id
		LEFT JOIN media am on am.id = ans.image
		WHERE q.id = $1
	`

	var q question.Detailed

	err := r.Pool.QueryRow(ctx, sql, questionId).Scan(
		&q.Text,
		&q.Answer,
		&q.AnswerMediaURL,
		&q.Author,
		&q.MediaURL,
		&q.MediaType,
		&q.LanguageId,
		&q.CreatedAt,
	)
	if err != nil {

		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("psql - question - FindById - r.c.Pool.QueryRow.Scan: %w", question.ErrNotFound)
		}

		return nil, fmt.Errorf("psql - question - FindById - r.c.Pool.QueryRow.Scan: %w", err)
	}

	q.Id = uint32(questionId)

	return &q, nil
}
