package question

import (
	"context"
	"errors"

	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) GetQuestion(ctx context.Context, p *pb.GetQuestionRequest) (*pb.GetQuestionResponse, error) {
	if p.QuestionId == 0 {
		return nil, twirp.RequiredArgumentError("question_id")
	}

	q, err := h.question.GetOne(ctx, p.QuestionId)
	if err != nil {
		if errors.Is(err, apperr.QuestionNotFound) {
			return nil, twirp.NotFoundError(err.Error())
		}

		return nil, twirp.InternalError(err.Error())
	}

	return &pb.GetQuestionResponse{
		Question: &pb.Question{
			Id:         q.ID,
			Text:       q.Text,
			Author:     q.Author,
			MediaUrl:   q.MediaURL,
			CreateTime: timestamppb.New(q.CreateTime),
			Answer: &pb.Answer{
				Id:       q.Answer.ID,
				Text:     q.Answer.Text,
				MediaUrl: q.Answer.MediaURL,
			},
		},
	}, nil
}
