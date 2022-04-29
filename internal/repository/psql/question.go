package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/answersuck/vault/internal/domain"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

const (
	questionTable      = "question"
	questionMediaTable = "question_media"
	answerTable        = "answer"
	answerImageTable   = "answer_image"
)

type question struct {
	log    logging.Logger
	client *postgres.Client
}

func NewQuestion(l logging.Logger, c *postgres.Client) *question {
	return &question{
		log:    l,
		client: c,
	}
}

func (r *question) FindAll(ctx context.Context) ([]*domain.Question, error) {
	sql := fmt.Sprintf(`
		SELECT
			q.id,
			q.question,
			a.answer,
			ai.url,
			acc.username,
			qm.url,
			qm.type,
			q.language_id,
			q.created_at
		FROM %s q
		LEFT JOIN %s a ON a.id = q.answer_id
		LEFT JOIN %s ai ON ai.id = a.answer_image_id
		LEFT JOIN %s acc ON acc.id = q.account_id
		LEFT JOIN %s qm ON qm.id = q.media_id
	`, questionTable, answerTable, answerImageTable, accountTable, questionMediaTable)

	r.log.Info("psql - question - FindAll: %s", sql)

	rows, err := r.client.Pool.Query(ctx, sql)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("r.Pool.Query: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("r.Pool.QueryRow.Scan: %w", err)
	}

	defer rows.Close()

	var questions []*domain.Question

	for rows.Next() {
		var q domain.Question

		if err = rows.Scan(
			&q.Id,
			&q.Q,
			&q.Answer,
			&q.AnswerImage,
			&q.Author,
			&q.Media,
			&q.MediaType,
			&q.LanguageId,
			&q.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		questions = append(questions, &q)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return questions, nil
}
