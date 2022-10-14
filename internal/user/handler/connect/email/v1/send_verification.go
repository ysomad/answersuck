package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	v1 "github.com/ysomad/answersuck/rpc/user/email/v1"
)

func (s *server) SendVerification(context.Context, *connect.Request[v1.SendVerificationRequest]) (*connect.Response[v1.SendVerificationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.email.v1.EmailService.SendVerification is not implemented"))
}
