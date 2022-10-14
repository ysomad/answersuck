package email

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	"github.com/ysomad/answersuck/user/internal/gen/proto/email"
)

func (s *server) VerifyEmail(ctx context.Context, r *connect.Request[email.VerifyEmailRequest]) (*connect.Response[email.VerifyEmailResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.email.EmailService.VerifyEmail is not implemented"))
}
