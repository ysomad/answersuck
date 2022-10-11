package accountv1

import (
	"context"
	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck-backend/api/proto/account/v1"
	"time"
)

func (s *server) CreateAccount(ctx context.Context, r *pb.CreateAccountRequest) (*pb.Account, error) {
	if r.Nickname == "" {
		return nil, twirp.RequiredArgumentError("nickname")
	}
	if r.Email == "" {
		return nil, twirp.RequiredArgumentError("email")
	}
	if r.Password == "" {
		return nil, twirp.RequiredArgumentError("password")
	}

	a, err := s.account.Create(ctx, r.Email, r.Nickname, r.Password)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	return &pb.Account{
		Id:         a.ID,
		Email:      a.Email,
		Nickname:   a.Nickname,
		IsVerified: a.Verified,
		IsArchived: a.Archived,
		CreatedAt:  a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  a.CreatedAt.Format(time.RFC3339),
	}, nil
}
