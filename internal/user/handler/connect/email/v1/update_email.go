package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	v1 "github.com/ysomad/answersuck/rpc/user/email/v1"
)

func (s *server) UpdateEmail(context.Context, *connect.Request[v1.UpdateEmailRequest]) (*connect.Response[v1.UpdateEmailResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.email.v1.EmailService.UpdateEmail is not implemented"))
}
