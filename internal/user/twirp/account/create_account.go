package account

import (
	"context"
	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck/rpc/user"
)

func (s *server) CreateAccount(ctx context.Context, r *pb.CreateAccountRequest) (*pb.Account, error) {
	if err := r.Validate(); err != nil {
		if pberr, ok := err.(pb.CreateAccountRequestValidationError); ok {
			return nil, twirp.InvalidArgumentError(pberr.Field(), pberr.Error())
		}

		return nil, twirp.NewError(twirp.InvalidArgument, err.Error())
	}

	return &pb.Account{}, nil
}
