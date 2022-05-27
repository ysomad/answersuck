package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/answersuck/vault/internal/domain/question"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
)

const (
	questionTable = "question"
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

func (r *questionPSQL) Save(ctx context.Context, dto *question.CreateDTO) (*question.Question, error) {
	sql := fmt.Sprintf(`
		WITH q AS (
			INSERT INTO %s(
				 text, 
				 answer_id, 
				 account_id, 
				 media_id, 
				 language_id, 
				 created_at, 
				 updated_at
			)
			VALUES (
				 $1::varchar, 
				 $2::integer, 
				 $3::uuid, 
				 $4::uuid, 
				 $5::integer, 
				 $6::timestamptz, 
				 $7::timestamptz
			)
			RETURNING id AS question_id, media_id, account_id
		)
		SELECT
			q.question_id,
			m.url AS media_url,
			m.mime_type AS media_type,
			a.username AS author
		FROM q
		LEFT JOIN media m ON m.id = q.media_id
		INNER JOIN account a ON a.id = q.account_id
	`, questionTable)

	r.log.Info("psql - question - Save: %s", sql)

	var q question.Question

	err := r.client.Pool.QueryRow(
		ctx,
		sql,
		dto.Text,
		dto.AnswerId,
		dto.AccountId,
		dto.MediaId,
		dto.LanguageId,
		dto.CreatedAt,
		dto.UpdatedAt,
	).Scan(&q.Id, &q.MediaURL, &q.MediaType, &q.Author)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return nil, fmt.Errorf("psql - r.client.Pool.QueryRow: %w", question.ErrForeignKeyViolation)
			}
		}

		return nil, fmt.Errorf("psql - r.client.Pool.QueryRow: %w", err)
	}

	q.Text = dto.Text
	q.AnswerId = dto.AnswerId
	q.LanguageId = dto.LanguageId
	q.CreatedAt = dto.CreatedAt
	q.UpdatedAt = dto.UpdatedAt

	return &q, nil
}

func (r *questionPSQL) FindAll(ctx context.Context) ([]*question.Question, error) {
	sql := fmt.Sprintf(`
	
	`, questionTable)

	r.log.Info("psql - question - FindAll: %s", sql)

	rows, err := r.client.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("r.client.Pool.Query.Scan: %w", err)
	}

	defer rows.Close()

	var questions []*question.Question

	for rows.Next() {
		var q question.Question

		if err = rows.Scan(); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		questions = append(questions, &q)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return questions, nil
}
