package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

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

func (r *questionPSQL) Save(ctx context.Context, q *question.Question) (int, error) {
	sql := fmt.Sprintf(`
		 INSERT INTO %s(
			  text, 
			  answer_id, 
			  account_id, 
			  media_id, 
			  language_id, 
			  created_at, 
			  updated_at
		 )
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id
	`, questionTable)

	r.log.Info("psql - question - Save: %s", sql)

	var questionId int

	err := r.client.Pool.QueryRow(
		ctx,
		sql,
		q.Text,
		q.AnswerId,
		q.AccountId,
		q.MediaId,
		q.LanguageId,
		q.CreatedAt,
		q.UpdatedAt,
	).Scan(&questionId)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return 0, fmt.Errorf("psql - r.client.Pool.QueryRow: %w", question.ErrForeignKeyViolation)
			}
		}

		return 0, fmt.Errorf("psql - r.client.Pool.QueryRow: %w", err)
	}

	return questionId, nil
}

func (r *questionPSQL) FindAll(ctx context.Context) ([]question.Minimized, error) {
	sql := fmt.Sprintf(`SELECT q.id, q.text, q.language_id FROM %s q`, questionTable)

	r.log.Info("psql - question - FindAll: %s", sql)

	rows, err := r.client.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("psql - r.client.Pool.Query.Scan: %w", err)
	}

	defer rows.Close()

	var qs []question.Minimized

	for rows.Next() {
		var q question.Minimized

		if err = rows.Scan(&q.Id, &q.Text, &q.LanguageId); err != nil {
			return nil, fmt.Errorf("psql - rows.Scan: %w", err)
		}

		qs = append(qs, q)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("psql - rows.Err: %w", err)
	}

	return qs, nil
}

func (r *questionPSQL) FindById(ctx context.Context, questionId int) (*question.Detailed, error) {
	sql := fmt.Sprintf(`
		SELECT
			q.text,
			ans.text AS answer,
			am.url AS answer_image_url,
			acc.username AS author,
			qm.url AS media_url,
			qm.mime_type AS media_type,
			q.language_id,
			q.created_at,
			q.updated_at
		FROM %s q
		INNER JOIN %s acc on acc.id = q.account_id
		INNER JOIN %s ans on ans.id = q.answer_id
		LEFT JOIN %s qm on qm.id = q.media_id
		LEFT JOIN %s am on am.id = ans.image
		WHERE q.id = $1
	`, questionTable, accountTable, answerTable, mediaTable, mediaTable)

	r.log.Info("psql - question - FindById: %s", sql)

	q := question.Detailed{Id: questionId}

	err := r.client.Pool.QueryRow(ctx, sql, questionId).Scan(
		&q.Text,
		&q.Answer,
		&q.AnswerImageURL,
		&q.Author,
		&q.MediaURL,
		&q.MediaType,
		&q.LanguageId,
		&q.CreatedAt,
		&q.UpdatedAt,
	)
	if err != nil {

		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("psql - r.client.Pool.QueryRow.Scan: %w", question.ErrNotFound)
		}

		return nil, fmt.Errorf("psql - r.client.Pool.QueryRow.Scan: %w", err)
	}

	return &q, nil
}
