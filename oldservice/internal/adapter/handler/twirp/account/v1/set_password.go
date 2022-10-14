package accountv1

import (
	"context"
	pb "github.com/ysomad/answersuck-backend/api/proto/account/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) SetPassword(ctx context.Context, r *pb.SetPasswordRequest) (*emptypb.Empty, error) {
	return nil, nil
}
