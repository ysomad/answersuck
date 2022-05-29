package psql

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

type questionRepo struct {
	l logging.Logger
	c *postgres.Client
}

func NewQuestionRepo(l logging.Logger, c *postgres.Client) *questionRepo {
	return &questionRepo{
		l: l,
		c: c,
	}
}

func (r *questionRepo) Save(ctx context.Context, q *question.Question) (int, error) {
	sql := `
		 INSERT INTO question(
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
	`

	r.l.Info("psql - question - Save: %s", sql)

	var questionId int

	err := r.c.Pool.QueryRow(
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
				return 0, fmt.Errorf("psql - question - Save - r.c.Pool.QueryRow.Scan: %w", question.ErrForeignKeyViolation)
			}
		}

		return 0, fmt.Errorf("psql - question - Save - r.c.Pool.QueryRow.Scan: %w", err)
	}

	return questionId, nil
}

func (r *questionRepo) FindAll(ctx context.Context) ([]question.Minimized, error) {
	sql := "SELECT id, text, language_id FROM question"

	r.l.Info("psql - question - FindAll: %s", sql)

	rows, err := r.c.Pool.Query(ctx, sql)
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

func (r *questionRepo) FindById(ctx context.Context, questionId int) (*question.Detailed, error) {
	sql := `
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
		FROM question q
		INNER JOIN account acc on acc.id = q.account_id
		INNER JOIN answer ans on ans.id = q.answer_id
		LEFT JOIN media qm on qm.id = q.media_id
		LEFT JOIN media am on am.id = ans.image
		WHERE q.id = $1
	`

	r.l.Info("psql - question - FindById: %s", sql)

	var q question.Detailed

	err := r.c.Pool.QueryRow(ctx, sql, questionId).Scan(
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
			return nil, fmt.Errorf("psql - question - FindById - r.c.Pool.QueryRow.Scan: %w", question.ErrNotFound)
		}

		return nil, fmt.Errorf("psql - question - FindById - r.c.Pool.QueryRow.Scan: %w", err)
	}

	q.Id = questionId

	return &q, nil
}
