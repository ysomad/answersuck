package email

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	"github.com/ysomad/answersuck/user/internal/gen/proto/email"
)

func (s *server) UpdateEmail(ctx context.Context, r *connect.Request[email.UpdateEmailRequest]) (*connect.Response[email.UpdateEmailResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.email.EmailService.UpdateEmail is not implemented"))
}
