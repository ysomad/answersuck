package v1

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"github.com/ysomad/answersuck/internal/pkg/session"
	apptwirp "github.com/ysomad/answersuck/internal/twirp"
	"github.com/ysomad/answersuck/internal/twirp/common"
	"github.com/ysomad/answersuck/internal/twirp/hooks"
	"github.com/ysomad/answersuck/internal/twirp/middleware"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	_ apptwirp.Handler   = &QuestionHandler{}
	_ pb.QuestionService = &QuestionHandler{}
)

type QuestionUseCase interface {
	Save(context.Context, *entity.Question) (int32, error)
	GetOne(context.Context, int32) (*entity.Question, error)
}

type QuestionHandler struct {
	question QuestionUseCase
	session  *session.Manager
}

func NewQuestionHandler(
	uc QuestionUseCase, sm *session.Manager) *QuestionHandler {
	return &QuestionHandler{
		question: uc,
		session:  sm,
	}
}

func (h *QuestionHandler) Handle(m *http.ServeMux) {
	s := pb.NewQuestionServiceServer(h,
		twirp.WithServerHooks(hooks.WithSession(h.session)))
	m.Handle(s.PathPrefix(), middleware.WithSessionID(s))
}

func (h *QuestionHandler) CreateQuestion(
	ctx context.Context,
	r *pb.CreateQuestionRequest) (*pb.CreateQuestionResponse, error) {
	session, err := common.CheckPlayerVerification(ctx)
	if err != nil {
		return nil, err
	}

	if r.Question == "" {
		return nil, twirp.RequiredArgumentError("question")
	}

	if r.Answer == "" {
		return nil, twirp.RequiredArgumentError("answer")
	}

	if err := r.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	questionID, err := h.question.Save(ctx, &entity.Question{
		Text:       r.Question,
		Author:     session.User.ID,
		MediaURL:   r.QuestionMediaUrl,
		CreateTime: time.Now(),
		Answer: entity.Answer{
			Text:     r.Answer,
			MediaURL: r.AnswerMediaUrl,
		},
	})
	if err != nil {
		if errors.Is(err, apperr.MediaNotFound) {
			return nil, twirp.InvalidArgument.Error(apperr.MsgQuestionMediaNotFound)
		}

		return nil, twirp.InternalError(err.Error())
	}

	return &pb.CreateQuestionResponse{QuestionId: questionID}, nil
}

func (h *QuestionHandler) GetQuestion(
	ctx context.Context,
	r *pb.GetQuestionRequest) (*pb.GetQuestionResponse, error) {
	if r.QuestionId == 0 {
		return nil, twirp.RequiredArgumentError("question_id")
	}

	q, err := h.question.GetOne(ctx, r.QuestionId)
	if err != nil {
		if errors.Is(err, apperr.QuestionNotFound) {
			return nil, twirp.NotFoundError(err.Error())
		}

		return nil, twirp.InternalError(err.Error())
	}

	return &pb.GetQuestionResponse{
		Question: &pb.Question{
			Id:         q.ID,
			Text:       q.Text,
			Author:     q.Author,
			MediaUrl:   q.MediaURL,
			CreateTime: timestamppb.New(q.CreateTime),
			Answer: &pb.Answer{
				Id:       q.Answer.ID,
				Text:     q.Answer.Text,
				MediaUrl: q.Answer.MediaURL,
			},
		},
	}, nil
}
