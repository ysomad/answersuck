package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"github.com/ysomad/answersuck/internal/pkg/session"
	apptwirp "github.com/ysomad/answersuck/internal/twirp"
	"github.com/ysomad/answersuck/internal/twirp/common"
	"github.com/ysomad/answersuck/internal/twirp/hooks"
	"github.com/ysomad/answersuck/internal/twirp/middleware"
)

var (
	_ apptwirp.Handler = &RoundHandler{}
	_ pb.RoundService  = &RoundHandler{}
)

type RoundUseCase interface {
	Create(ctx context.Context, r entity.Round) (roundID int32, err error)
	Update(ctx context.Context, r entity.Round) error
	GetAll(ctx context.Context, packID int32) ([]entity.Round, error)
	AddTopic(ctx context.Context, roundID, topicID int32) (int32, error)
	RemoveTopic(ctx context.Context, roundID, topicID int32) error
}

type RoundHandler struct {
	round   RoundUseCase
	session *session.Manager
}

func NewRoundHandler(uc RoundUseCase, sm *session.Manager) *RoundHandler {
	return &RoundHandler{
		round:   uc,
		session: sm,
	}
}

func (h *RoundHandler) Handle(m *http.ServeMux) {
	s := pb.NewRoundServiceServer(h,
		twirp.WithServerHooks(hooks.WithSession(h.session)))
	m.Handle(s.PathPrefix(), middleware.WithSessionID(s))
}

func (h *RoundHandler) CreateRound(
	ctx context.Context,
	r *pb.CreateRoundRequest) (*pb.CreateRoundResponse, error) {
	var err error

	if _, err = common.CheckPlayerVerification(ctx); err != nil {
		return nil, err
	}

	if r.PackId == 0 {
		return nil, twirp.RequiredArgumentError("pack_id")
	}

	if r.RoundName == "" {
		return nil, twirp.RequiredArgumentError("round_name")
	}

	if r.RoundPosition == 0 {
		return nil, twirp.RequiredArgumentError("round_position")
	}

	if err := r.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	round := entity.Round{
		Name:     r.RoundName,
		PackID:   r.PackId,
		Position: int16(r.RoundPosition),
	}

	round.ID, err = h.round.Create(ctx, round)
	if err != nil {
		switch {
		case errors.Is(err, apperr.PackNotFound):
			return nil, twirp.InvalidArgumentError("pack_id", apperr.MsgPackNotFound)
		case errors.Is(err, apperr.PackNotAuthor):
			return nil, twirp.PermissionDenied.Error(apperr.MsgPackNotAuthor)
		}

		return nil, twirp.InternalError(err.Error())
	}

	return &pb.CreateRoundResponse{
		Round: &pb.Round{
			Id:       round.ID,
			Name:     round.Name,
			Position: int32(round.Position),
			PackId:   round.PackID,
		},
	}, nil
}

func (h *RoundHandler) UpdateRound(
	ctx context.Context,
	r *pb.UpdateRoundRequest) (*pb.UpdateRoundResponse, error) {
	if _, err := common.CheckPlayerVerification(ctx); err != nil {
		return nil, err
	}

	if r.RoundId == 0 {
		return nil, twirp.RequiredArgumentError("round_id")
	}

	if r.RoundName == "" {
		return nil, twirp.RequiredArgumentError("round_name")
	}

	if r.RoundPosition == 0 {
		return nil, twirp.RequiredArgumentError("round_position")
	}

	if r.PackId == 0 {
		return nil, twirp.RequiredArgumentError("pack_id")
	}

	if err := r.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	round := entity.Round{
		ID:       r.RoundId,
		Name:     r.RoundName,
		PackID:   r.PackId,
		Position: int16(r.RoundPosition),
	}

	if err := h.round.Update(ctx, round); err != nil {
		switch {
		case errors.Is(err, apperr.PackNotAuthor):
			return nil, twirp.PermissionDenied.Error(apperr.MsgPackNotAuthor)
		case errors.Is(err, apperr.RoundNotFound):
			return nil, twirp.InvalidArgument.Error(apperr.MsgRoundNotFound)
		case errors.Is(err, apperr.PackNotFound):
			return nil, twirp.InvalidArgument.Error(apperr.MsgPackNotFound)
		}

		return nil, twirp.InternalError(err.Error())
	}

	return &pb.UpdateRoundResponse{
		Round: &pb.Round{
			Id:       round.ID,
			Name:     round.Name,
			Position: int32(round.Position),
			PackId:   round.PackID,
		},
	}, nil
}

func (h *RoundHandler) ListRounds(
	ctx context.Context,
	r *pb.ListRoundsRequest) (*pb.ListRoundsResponse, error) {
	if r.PackId == 0 {
		return nil, twirp.RequiredArgumentError("pack_id")
	}

	rr, err := h.round.GetAll(ctx, r.PackId)
	if err != nil {
		if errors.Is(err, apperr.PackNotFound) {
			return nil, twirp.InvalidArgument.Error(apperr.MsgPackNotFound)
		}

		return nil, twirp.InternalError(err.Error())
	}

	rounds := make([]*pb.Round, len(rr))

	for i, r := range rr {
		rounds[i] = &pb.Round{
			Id:       r.ID,
			Name:     r.Name,
			Position: int32(r.Position),
			PackId:   r.PackID,
		}
	}

	return &pb.ListRoundsResponse{Rounds: rounds}, nil
}

func (h *RoundHandler) AddTopic(
	ctx context.Context,
	r *pb.AddTopicRequest) (*pb.AddTopicResponse, error) {
	if _, err := common.CheckPlayerVerification(ctx); err != nil {
		return nil, err
	}

	if r.RoundId == 0 {
		return nil, twirp.RequiredArgumentError("round_id")
	}

	if r.TopicId == 0 {
		return nil, twirp.RequiredArgumentError("topic_id")
	}

	roundTopicID, err := h.round.AddTopic(ctx, r.RoundId, r.TopicId)
	if err != nil {
		switch {
		case errors.Is(err, apperr.RoundNotFound):
			return nil, twirp.InvalidArgument.Error(apperr.MsgRoundNotFound)
		case errors.Is(err, apperr.PackNotAuthor):
			return nil, twirp.PermissionDenied.Error(apperr.MsgPackNotAuthor)
		case errors.Is(err, apperr.RoundTopicNotAdded):
			return nil, twirp.InvalidArgument.Error(apperr.MsgRoundTopicNotAdded)
		case errors.Is(err, apperr.RoundTopicAlreadyExists):
			return nil, twirp.InvalidArgument.Error(apperr.MsgRoundTopicAlreadyExists)
		}

		return nil, twirp.InternalError(err.Error())
	}

	return &pb.AddTopicResponse{
		RoundTopicId: roundTopicID,
	}, nil
}

func (h *RoundHandler) RemoveTopic(
	ctx context.Context,
	r *pb.RemoveTopicRequest) (*emptypb.Empty, error) {
	if _, err := common.CheckPlayerVerification(ctx); err != nil {
		return nil, err
	}

	if r.RoundId == 0 {
		return nil, twirp.RequiredArgumentError("round_id")
	}

	if r.TopicId == 0 {
		return nil, twirp.RequiredArgumentError("topic_id")
	}

	if err := h.round.RemoveTopic(ctx, r.RoundId, r.TopicId); err != nil {
		switch {
		case errors.Is(err, apperr.PackNotAuthor):
			return nil, twirp.PermissionDenied.Error(apperr.MsgPackNotAuthor)
		case errors.Is(err, apperr.RoundTopicNotDeleted):
			return nil, twirp.InvalidArgument.Error(apperr.MsgRoundTopicNotDeleted)
		}

		return nil, twirp.InternalError(err.Error())
	}

	return new(emptypb.Empty), nil
}

func (h *RoundHandler) GetQuestionGrid(
	ctx context.Context,
	r *pb.GetQuestionGridRequest) (*pb.GetQuestionGridResponse, error) {
	// TODO: IMPLEMENT
	return nil, nil
}
