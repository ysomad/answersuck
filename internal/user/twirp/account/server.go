package account

import (
	"context"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"

	pb "github.com/ysomad/answersuck/rpc/user"
)

var _ pb.AccountService = &server{}

type server struct{}

func NewServer() *server {
	return &server{}
}

func (s *server) CreateAccount(ctx context.Context, r *pb.CreateAccountRequest) (*pb.Account, error) {
	if err := r.ValidateAll(); err != nil {

		pberrs, ok := err.(pb.CreateAccountRequestMultiError)
		if ok {
			twirperr := twirp.NewError(twirp.InvalidArgument, "validation error")

			for _, err := range pberrs {
				pberr, ok := err.(pb.CreateAccountRequestValidationError)
				if ok {
					twirperr = twirperr.WithMeta(strings.ToLower(pberr.Field()), pberr.Error())
				}
			}

			return nil, twirperr
		}

		return nil, twirp.NewError(twirp.InvalidArgument, err.Error())
	}

	return &pb.Account{}, nil
}

func (s *server) GetAccountById(ctx context.Context, r *pb.GetAccountByIdRequest) (*pb.Account, error) {
	return &pb.Account{}, nil
}

func (s *server) GetAccountByEmail(ctx context.Context, r *pb.GetAccountByEmailRequest) (*pb.Account, error) {
	return &pb.Account{}, nil
}

func (s *server) DeleteAccount(ctx context.Context, r *pb.DeleteAccountRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
