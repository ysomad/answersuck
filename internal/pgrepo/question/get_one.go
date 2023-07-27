package question

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"golang.org/x/exp/slog"
)

func (r *Repository) GetOne(ctx context.Context, questionID int32) (*entity.Question, error) {
	sql, args, err := r.Builder.
		Select(
			"q.id as id",
			"q.text as text",
			"q.author as author",
			"q.media_url as media_url",
			"q.create_time as create_time",
			"a.id as answer_id",
			"a.text as answer",
			"a.media_url as answer_media_url").
		From(questionTable + " q").
		InnerJoin(answerTable + " a ON q.answer_id = a.id").
		Where(squirrel.Eq{"q.id": questionID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	slog.Info(sql)

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("r.Pool.Query: %w", err)
	}

	q, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[question])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.QuestionNotFound
		}

		return nil, fmt.Errorf("pgx.CollectOneRow: %w", err)
	}

	return &entity.Question{
		ID:   q.ID,
		Text: q.Text,
		Answer: entity.Answer{
			ID:       q.AnswerID,
			Text:     q.Answer,
			MediaURL: string(q.AnswerMediaURL),
		},
		Author:     q.Author,
		MediaURL:   string(q.MediaURL),
		CreateTime: q.CreateTime,
	}, nil
}
