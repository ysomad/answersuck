package password

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ysomad/answersuck/user/internal/gen/proto/password"
)

func (s *server) ResetPassword(ctx context.Context, r *connect.Request[password.ResetPasswordRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.password.PasswordService.ResetPassword is not implemented"))
}
