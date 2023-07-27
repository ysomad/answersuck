package question

import (
	"context"
	"errors"
	"time"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"github.com/ysomad/answersuck/internal/twirp/common"
)

func (h *Handler) CreateQuestion(ctx context.Context, p *pb.CreateQuestionRequest) (*pb.CreateQuestionResponse, error) {
	session, err := common.CheckPlayerVerification(ctx)
	if err != nil {
		return nil, err
	}

	if p.Question == "" {
		return nil, twirp.RequiredArgumentError("question")
	}

	if p.Answer == "" {
		return nil, twirp.RequiredArgumentError("answer")
	}

	if err := p.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	questionID, err := h.question.Save(ctx, &entity.Question{
		Text: p.Question,
		Answer: entity.Answer{
			Text:  p.Answer,
			Media: entity.Media{URL: p.AnswerMediaUrl},
		},
		Author:     session.User.UserID,
		Media:      entity.Media{URL: p.QuestionMediaUrl},
		CreateTime: time.Now(),
	})
	if err != nil {
		switch {
		case errors.Is(err, apperr.ErrAnswerMediaNotExist):
			return nil, twirp.InvalidArgumentError("answer_media_url", err.Error())
		case errors.Is(err, apperr.ErrQestionMediaNotExist):
			return nil, twirp.InvalidArgumentError("question_media_url", err.Error())
		}

		return nil, twirp.InternalError(err.Error())
	}

	return &pb.CreateQuestionResponse{QuestionId: questionID}, nil
}
