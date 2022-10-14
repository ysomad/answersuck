package accountv1

import (
	"context"
	pb "github.com/ysomad/answersuck-backend/api/proto/account/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) ResetPassword(ctx context.Context, r *pb.ResetPasswordRequest) (*emptypb.Empty, error) {
	return nil, nil
}
