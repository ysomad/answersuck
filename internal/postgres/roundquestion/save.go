package roundquestion

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype/zeronull"
	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (r *Repository) Save(ctx context.Context, q *entity.RoundQuestion) (int32, error) {
	sql, args, err := r.Builder.
		Insert(RoundQuestionsTable).
		Columns(
			"round_id",
			"topic_id",
			"question_id",
			"question_type",
			"cost",
			"answer_time",
			"host_comment",
			"secret_topic",
			"secret_cost",
			"transfer_type",
			"is_keepable",
		).
		Values(
			q.RoundID,
			q.TopicID,
			q.QuestionID,
			q.Type,
			q.Cost,
			q.AnswerTime,
			zeronull.Text(q.HostComment),
			zeronull.Text(q.SecretTopic),
			zeronull.Int4(q.SecretCost),
			zeronull.Int2(q.TransferType),
			q.Keepable,
		).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, err
	}

	var id int32

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "round_questions_round_id_fkey":
				return 0, apperr.RoundNotFound
			case "round_questions_topic_id_fkey":
				return 0, apperr.TopicNotFound
			case "round_questions_question_id_fkey":
				return 0, apperr.QuestionNotFound
			}
		}

		return 0, fmt.Errorf("error saving round question: %w", err)
	}

	return id, nil
}
