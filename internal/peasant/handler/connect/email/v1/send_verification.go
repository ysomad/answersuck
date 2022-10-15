package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) SendVerification(context.Context, *connect.Request[pb.SendVerificationRequest]) (*connect.Response[pb.SendVerificationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.email.v1.EmailService.SendVerification is not implemented"))
}
