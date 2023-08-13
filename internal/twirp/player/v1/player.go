package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/player/v1"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	apptwirp "github.com/ysomad/answersuck/internal/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	_ apptwirp.Handler = &PlayerHandler{}
	_ pb.PlayerService = &PlayerHandler{}
)

type PlayerUseCase interface {
	Create(ctx context.Context, nickname, email, password string) error
	Get(ctx context.Context, login string) (*entity.Player, error)
}

type PlayerHandler struct {
	player PlayerUseCase
}

func NewPlayerHandler(uc PlayerUseCase) *PlayerHandler {
	return &PlayerHandler{
		player: uc,
	}
}

func (h *PlayerHandler) Handle(m *http.ServeMux) {
	s := pb.NewPlayerServiceServer(h)
	m.Handle(s.PathPrefix(), s)
}

func (h *PlayerHandler) GetPlayer(
	ctx context.Context, r *pb.GetPlayerRequest) (*pb.GetPlayerResponse, error) {
	if r.Nickname == "" {
		return nil, twirp.RequiredArgumentError("nickname")
	}

	player, err := h.player.Get(ctx, r.Nickname)
	if err != nil {
		if errors.Is(err, apperr.PlayerNotFound) {
			return nil, twirp.NotFoundError(apperr.PlayerNotFound.Error())
		}

		return nil, twirp.InternalError(err.Error())
	}

	return &pb.GetPlayerResponse{
		Player: &pb.Player{
			Nickname:      player.Nickname,
			Email:         player.Email,
			DisplayName:   player.DisplayName,
			EmailVerified: player.EmailVerified,
			CreateTime:    timestamppb.New(player.CreatedAt),
		},
	}, nil
}

func (h *PlayerHandler) CreatePlayer(
	ctx context.Context, r *pb.CreatePlayerRequest) (*emptypb.Empty, error) {
	if r.Email == "" {
		return nil, twirp.RequiredArgumentError("email")
	}

	if r.Nickname == "" {
		return nil, twirp.RequiredArgumentError("nickname")
	}

	if r.Password == "" {
		return nil, twirp.RequiredArgumentError("password")
	}

	if err := r.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	if err := h.player.Create(ctx, r.Nickname, r.Email, r.Password); err != nil {
		if errors.Is(err, apperr.PlayerAlreadyExists) {
			return nil, twirp.AlreadyExists.Error(err.Error())
		}

		return nil, twirp.InternalError(err.Error())
	}

	return &emptypb.Empty{}, nil
}
