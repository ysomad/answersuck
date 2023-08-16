package roundquestion

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (r *Repository) GetOne(
	ctx context.Context, id int32) (*entity.RoundQuestionDetailed, error) {
	sql, args, err := r.Builder.
		Select(
			"rq.id as id",
			"rq.question_type as question_type",
			"rq.cost as question_cost",
			"rq.answer_time as answer_time",
			"rq.host_comment as host_comment",
			"rq.secret_topic as secret_topic",
			"rq.secret_cost as secret_cost",
			"rq.is_keepable as is_keepable",
			"rq.transfer_type as transfer_type",
			"rq.question_id as question_id",
			"q.text as question",
			"q.media_url as question_media_url",
			"q.answer_id as answer_id",
			"a.text as answer",
			"a.media_url as answer_media_url",
			"rp.round_id as round_id",
			"rp.topic_id as topic_id").
		From("round_questions rq").
		InnerJoin("questions q ON rq.question_id = q.id").
		InnerJoin("answers a ON q.answer_id = a.id").
		InnerJoin("round_topics rp ON rq.round_topic_id = rp.id").
		Where(squirrel.Eq{"rq.id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	q, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[roundQuestion])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.RoundQuestionNotFound
		}

		return nil, err
	}

	return &entity.RoundQuestionDetailed{
		RoundQuestion: entity.RoundQuestion{
			ID:           q.ID,
			QuestionID:   q.QuestionID,
			TopicID:      q.TopicID,
			RoundID:      q.RoundID,
			Type:         q.Type,
			Cost:         q.Cost,
			AnswerTime:   q.AnswerTime,
			HostComment:  string(q.HostComment),
			SecretTopic:  string(q.SecretTopic),
			SecretCost:   int32(q.SecretCost),
			Keepable:     q.Keepable.Bool,
			TransferType: entity.QuestionTransferType(q.TransferType),
		},
		Question:         q.Question,
		QuestionMediaURL: string(q.QuestionMediaURL),
		AnswerID:         q.AnswerID,
		Answer:           q.Answer,
		AnswerMediaURL:   string(q.AnswerMediaURL),
	}, nil
}
