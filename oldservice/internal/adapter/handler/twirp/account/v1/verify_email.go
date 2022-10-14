package accountv1

import (
	"context"
	pb "github.com/ysomad/answersuck-backend/api/proto/account/v1"
)

func (s *server) VerifyEmail(ctx context.Context, r *pb.VerifyEmailRequest) (*pb.Account, error) {
	return nil, nil
}
