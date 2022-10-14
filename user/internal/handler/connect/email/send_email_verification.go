package email

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ysomad/answersuck/user/internal/gen/proto/email"
)

func (s *server) SendEmailVerification(ctx context.Context, r *connect.Request[email.SendEmailVerificationRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.email.EmailService.SendEmailVerification is not implemented"))
}
