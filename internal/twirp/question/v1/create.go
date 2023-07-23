package v1

import (
	"context"

	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck/internal/gen/api/question/v1"
	"github.com/ysomad/answersuck/internal/twirp/common"
)

func (h *Handler) CreateQuestion(ctx context.Context, p *pb.CreateQuestionRequest) (*pb.CreateQuestionResponse, error) {
	if _, err := common.CheckPlayerVerification(ctx); err != nil {
		return nil, err
	}

	if p.Question == "" {
		return nil, twirp.RequiredArgumentError("question")
	}

	if p.Answer.Text == "" {
		return nil, twirp.RequiredArgumentError("answer.text")
	}

	return nil, nil
}
