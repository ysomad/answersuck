package accountv1

import (
	"context"
	pb "github.com/ysomad/answersuck-backend/api/proto/account/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) RequestEmailVerification(ctx context.Context, r *pb.RequestEmailVerificationRequest) (*emptypb.Empty, error) {
	return nil, nil
}
