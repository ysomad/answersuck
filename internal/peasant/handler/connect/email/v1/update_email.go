package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) UpdateEmail(context.Context, *connect.Request[pb.UpdateEmailRequest]) (*connect.Response[pb.UpdateEmailResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.email.v1.EmailService.UpdateEmail is not implemented"))
}
