package v1

import (
	"context"
	"net/http"

	"github.com/twitchtv/twirp"

	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/pkg/session"
	apptwirp "github.com/ysomad/answersuck/internal/twirp"
	"github.com/ysomad/answersuck/internal/twirp/common"
	"github.com/ysomad/answersuck/internal/twirp/hooks"
	"github.com/ysomad/answersuck/internal/twirp/middleware"
)

var (
	_ apptwirp.Handler        = &RoundQuestionHandler{}
	_ pb.RoundQuestionService = &RoundQuestionHandler{}
)

type RoundQuestionUseCase interface {
	Save(ctx context.Context, q *entity.RoundQuestion) (int32, error)
}

type RoundQuestionHandler struct {
	round   RoundQuestionUseCase
	session *session.Manager
}

func NewRoundQuestionHandler(
	uc RoundQuestionUseCase, sm *session.Manager) *RoundQuestionHandler {
	return &RoundQuestionHandler{
		round:   uc,
		session: sm,
	}
}

func (h *RoundQuestionHandler) Handle(m *http.ServeMux) {
	s := pb.NewRoundQuestionServiceServer(h,
		twirp.WithServerHooks(hooks.WithSession(h.session)))
	m.Handle(s.PathPrefix(), middleware.WithSessionID(s))
}

func (h *RoundQuestionHandler) CreateRoundQuestion(ctx context.Context,
	r *pb.CreateRoundQuestionRequest) (*pb.CreateRoundQuestionResponse, error) {
	var err error

	if _, err = common.CheckPlayerVerification(ctx); err != nil {
		return nil, err
	}

	if r.QuestionId == 0 {
		return nil, twirp.RequiredArgumentError("question_id")
	}

	if r.TopicId == 0 {
		return nil, twirp.RequiredArgumentError("topic_id")
	}

	if r.RoundId == 0 {
		return nil, twirp.RequiredArgumentError("round_id")
	}

	if err = r.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	q := &entity.RoundQuestion{
		QuestionID:   r.QuestionId,
		TopicID:      r.TopicId,
		RoundID:      r.RoundId,
		Type:         entity.QuestionType(r.QuestionType),
		Cost:         r.QuestionCost,
		AnswerTime:   r.AnswerTime.AsDuration(),
		HostComment:  r.HostComment,
		SecretTopic:  r.SecretTopic,
		SecretCost:   r.SecretCost,
		Keepable:     r.IsKeepable,
		TransferType: entity.QuestionTransferType(r.TransferType),
	}

	if err = q.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	q.ID, err = h.round.Save(ctx, q)
	if err != nil {
		// TODO: handle specific errors
		return nil, twirp.InternalError(err.Error())
	}

	return &pb.CreateRoundQuestionResponse{RoundQuestionId: q.ID}, nil
}

func (h *RoundQuestionHandler) GetRoundQuestion(
	ctx context.Context,
	r *pb.GetRoundQuestionRequest) (*pb.GetRoundQuestionResponse, error) {
	// TODO: IMPLEMENT
	return nil, nil
}
