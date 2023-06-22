package playerv1

import (
	"context"

	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/ysomad/answersuck/internal/gen/proto/player/v1"
)

var _ pb.PasswordService = &passwordHandler{}

type passwordHandler struct{}

func NewPasswordHandler() *passwordHandler {
	return &passwordHandler{}
}

func (h *passwordHandler) ResetPassword(context.Context, *pb.ResetPasswordRequest) (*emptypb.Empty, error) {
	return new(emptypb.Empty), twirp.InternalError("not implemented")
}

func (h *passwordHandler) SetPassword(context.Context, *pb.SetPasswordRequest) (*emptypb.Empty, error) {
	return new(emptypb.Empty), twirp.InternalError("not implemented")
}

func (h *passwordHandler) UpdatePassword(context.Context, *pb.UpdatePasswordRequest) (*emptypb.Empty, error) {
	return new(emptypb.Empty), twirp.InternalError("not implemented")
}
