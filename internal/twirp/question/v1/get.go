package v1

import (
	"context"

	pb "github.com/ysomad/answersuck/internal/gen/api/question/v1"
)

func (h *Handler) GetQuestion(ctx context.Context, p *pb.GetQuestionRequest) (*pb.GetQuestionResponse, error) {
	return nil, nil
}
