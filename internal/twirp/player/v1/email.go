package playerv1

import (
	"context"

	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/ysomad/answersuck/internal/gen/proto/player/v1"
)

var _ pb.EmailService = &emailHandler{}

type emailHandler struct{}

func NewEmailHandler() *emailHandler {
	return &emailHandler{}
}

func (h *emailHandler) UpdateEmail(ctx context.Context, r *pb.UpdateEmailRequest) (*pb.UpdateEmailResponse, error) {
	return nil, twirp.InternalError("not implemented")
}

func (h *emailHandler) VerifyEmail(ctx context.Context, r *pb.VerifyEmailRequest) (*emptypb.Empty, error) {
	return new(emptypb.Empty), twirp.InternalError("not implemented")
}

func (h *emailHandler) SendVerification(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return new(emptypb.Empty), twirp.InternalError("not implemented")
}
