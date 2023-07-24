package question

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype/zeronull"
	"github.com/ysomad/answersuck/internal/entity"
)

func (r *Repository) Save(ctx context.Context, q *entity.Question) (questionID int32, err error) {
	txFunc := func(tx pgx.Tx) error {
		sql, args, err := r.Builder.
			Insert(answerTable).
			Columns("text, media_url").
			Values(q.Answer.Text, zeronull.Text(q.Answer.Media.URL)).
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
			Values(q.Text, q.Answer.ID, q.Author, zeronull.Text(q.Media.URL), q.CreateTime).
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
		return 0, err
	}

	return q.ID, nil
}
