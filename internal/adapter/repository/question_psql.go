package repository

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/domain/question"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

const (
	questionTable      = "question"
	questionMediaTable = "question_media"
	answerImageTable   = "answer_image"
)

type questionPSQL struct {
	log    logging.Logger
	client *postgres.Client
}

func NewQuestionPSQL(l logging.Logger, c *postgres.Client) *questionPSQL {
	return &questionPSQL{
		log:    l,
		client: c,
	}
}

func (r *questionPSQL) Create(ctx context.Context, dto *question.CreateDTO) (*question.Question, error) {
	_ = fmt.Sprintf(`
		INSERT INTO %s(question, answer_id, media_id, account_id, language_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, questionTable)

	return nil, nil
}

func (r *questionPSQL) FindAll(ctx context.Context) ([]*question.Question, error) {
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
		return nil, fmt.Errorf("r.client.Pool.Query.Scan: %w", err)
	}

	defer rows.Close()

	var questions []*question.Question

	for rows.Next() {
		var q question.Question

		if err = rows.Scan(
			&q.Id,
			&q.Text,
			&q.Answer,
			&q.AnswerImageURL,
			&q.Author,
			&q.MediaURL,
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
