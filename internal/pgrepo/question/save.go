package question

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype/zeronull"
	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (r *Repository) Save(ctx context.Context, q *entity.Question) (questionID int32, err error) {
	txFunc := func(tx pgx.Tx) error {
		sql, args, err := r.Builder.
			Insert(answerTable).
			Columns("text, media_url").
			Values(q.Answer.Text, zeronull.Text(q.Answer.MediaURL)).
			Suffix("RETURNING id").
			ToSql()
		if err != nil {
			return err
		}

		if err := tx.QueryRow(ctx, sql, args...).Scan(&q.Answer.ID); err != nil {
			return fmt.Errorf("error saving answer: %w", err)
		}

		sql, args, err = r.Builder.
			Insert(questionTable).
			Columns("text, answer_id, author, media_url, create_time").
			Values(q.Text, q.Answer.ID, q.Author, zeronull.Text(q.MediaURL), q.CreateTime).
			Suffix("RETURNING id").
			ToSql()
		if err != nil {
			return err
		}

		if err := tx.QueryRow(ctx, sql, args...).Scan(&q.ID); err != nil {
			return fmt.Errorf("error saving question: %w", err)
		}

		return nil
	}

	if err := pgx.BeginTxFunc(ctx, r.Pool, pgx.TxOptions{}, txFunc); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.ConstraintName == "questions_media_url_fkey" ||
			pgErr.ConstraintName == "answers_media_url_fkey" {
			return 0, apperr.MediaNotFound
		}

		return 0, err
	}

	return q.ID, nil
}
