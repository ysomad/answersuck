package repository

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/dto"

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

type questionRepo struct {
	log    logging.Logger
	client *postgres.Client
}

func NewQuestionRepo(l logging.Logger, c *postgres.Client) *questionRepo {
	return &questionRepo{
		log:    l,
		client: c,
	}
}

func (r *questionRepo) Create(ctx context.Context, qc *dto.QuestionCreate) (*domain.Question, error) {
	_ = fmt.Sprintf(`
		INSERT INTO %s(question, media_id, answer, language_id, account_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, questionTable)

	return nil, nil
}

func (r *questionRepo) FindAll(ctx context.Context) ([]*domain.Question, error) {
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

	var questions []*domain.Question

	for rows.Next() {
		var q domain.Question

		if err = rows.Scan(
			&q.Id,
			&q.Q,
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
