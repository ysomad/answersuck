package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) VerifyEmail(context.Context, *connect.Request[pb.VerifyEmailRequest]) (*connect.Response[pb.VerifyEmailResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.email.v1.EmailService.VerifyEmail is not implemented"))
}
