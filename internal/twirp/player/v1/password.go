package playerv1

import (
	"context"

	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/ysomad/answersuck/internal/gen/proto/player/v1"
)

var _ pb.PasswordService = &PasswordHandler{}

type PasswordHandler struct{}

func NewPasswordHandler() *PasswordHandler {
	return &PasswordHandler{}
}

func (h *PasswordHandler) ResetPassword(context.Context, *pb.ResetPasswordRequest) (*emptypb.Empty, error) {
	return new(emptypb.Empty), twirp.InternalError("not implemented")
}

func (h *PasswordHandler) SetPassword(context.Context, *pb.SetPasswordRequest) (*emptypb.Empty, error) {
	return new(emptypb.Empty), twirp.InternalError("not implemented")
}

func (h *PasswordHandler) UpdatePassword(context.Context, *pb.UpdatePasswordRequest) (*emptypb.Empty, error) {
	return new(emptypb.Empty), twirp.InternalError("not implemented")
}
