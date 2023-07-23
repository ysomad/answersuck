package question

import (
	"context"

	"github.com/ysomad/answersuck/internal/entity"
)

func (r *Repository) Save(ctx context.Context, q *entity.Question) (*entity.Question, error) {
	// txFunc := func(tx pgx.Tx) error {
	// 	// 1. Insert answer answerMedia
	// 	answerMedia, err := r.saveMedia(ctx, tx, q.Answer.Media)
	// 	if err != nil {
	// 		return fmt.Errorf("error saving answer media: %w", err)
	// 	}

	// 	// 2. Insert question media
	// 	questionMedia, err := r.saveMedia(ctx, tx, q.Media)
	// 	if err != nil {
	// 		return fmt.Errorf("error saving question media: %w", err)
	// 	}

	// 	// 3. Insert answer
	// 	answer, err := r.saveAnswer(ctx, tx, q.Answer)
	// 	if err != nil {
	// 		return fmt.Errorf("error saving question media: %w", err)
	// 	}

	// 	// 4. Insert question

	// 	return nil
	// }

	// if err := pgx.BeginTxFunc(ctx, r.Pool, pgx.TxOptions{}, txFunc); err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

// func (r *Repository) saveAnswer(ctx context.Context, tx pgx.Tx, a entity.Answer) (entity.Answer, error) {
// 	sql, args, err := r.Builder.
// 		Insert(answerTable).
// 		Columns("text, author, media_url, create_time").
// 		Values(a.Text, a.Author, a.Media.URL, a.CreateTime).
// 		Suffix("RETURNING id").
// 		ToSql()
// 	if err != nil {
// 		return entity.Answer{}, err
// 	}

// 	if err := tx.QueryRow(ctx, sql, args...).Scan(&a.ID); err != nil {
// 		return entity.Answer{}, fmt.Errorf("scan error: %w", err)
// 	}

// 	return a, nil
// }

// func (r *Repository) saveMedia(ctx context.Context, tx pgx.Tx, m entity.Media) (entity.Media, error) {
// 	sql, args, err := r.Builder.
// 		Insert(mediaTable).
// 		Columns("url, type, uploaded_by, create_time").
// 		Values(m.URL, m.Type, m.Author, m.CreateTime).
// 		Suffix("ON CONFLICT(url) DO UPDATE").
// 		Suffix("SET url = EXCLUDED.url").
// 		Suffix("RETURNING url, type, uploaded_by, create_time").
// 		ToSql()
// 	if err != nil {
// 		return entity.Media{}, err
// 	}

// 	rows, err := tx.Query(ctx, sql, args...)
// 	if err != nil {
// 		return entity.Media{}, err
// 	}

// 	res, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[media])
// 	if err != nil {
// 		return entity.Media{}, err
// 	}

// 	return entity.Media(res), nil
// }
