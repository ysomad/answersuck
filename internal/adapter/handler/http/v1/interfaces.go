package v1

import (
	"context"
	"io"

	"github.com/answersuck/host/internal/domain/account"
	"github.com/answersuck/host/internal/domain/answer"
	"github.com/answersuck/host/internal/domain/language"
	"github.com/answersuck/host/internal/domain/media"
	"github.com/answersuck/host/internal/domain/question"
	"github.com/answersuck/host/internal/domain/session"
	"github.com/answersuck/host/internal/domain/tag"
	"github.com/answersuck/host/internal/domain/topic"
	"github.com/answersuck/host/internal/pkg/pagination"
)

type validate interface {
	TranslateError(err error) map[string]string
	RequestBody(b io.ReadCloser, dest any) error
	Struct(s any) error
}

type accountService interface {
	Create(ctx context.Context, email, nickname, password string) (account.Account, error)
	Delete(ctx context.Context, accountId string) error
	RequestVerification(ctx context.Context, accountId string) error
	Verify(ctx context.Context, code string) error
	UpdatePassword(ctx context.Context, accountId, oldPwd, newPwd string) error
	ResetPassword(ctx context.Context, login string) error
	SetPassword(ctx context.Context, token, password string) error
}

type sessionService interface {
	GetByIdWithDetails(ctx context.Context, sessionId string) (*session.WithAccountDetails, error)
	GetById(ctx context.Context, sessionId string) (*session.Session, error)
	GetAll(ctx context.Context, accountId string) ([]*session.Session, error)
	Terminate(ctx context.Context, sessionId string) error
	TerminateAllWithExcept(ctx context.Context, accountId, sessionId string) error
}

type (
	loginService interface {
		Login(ctx context.Context, login, password string, d session.Device) (*session.Session, error)
	}

	tokenService interface {
		Create(ctx context.Context, accountId, password string) (string, error)
		Parse(ctx context.Context, token string) (string, error)
	}
)

type mediaService interface {
	UploadAndSave(ctx context.Context, m media.Media, size int64) (media.WithURL, error)
}

type languageService interface {
	GetAll(ctx context.Context) ([]language.Language, error)
}

type tagService interface {
	CreateMultiple(ctx context.Context, r []tag.Tag) ([]tag.Tag, error)
	GetAll(ctx context.Context, p tag.ListParams) (pagination.List[tag.Tag], error)
}

type topicService interface {
	Create(ctx context.Context, t topic.Topic) (topic.Topic, error)
	GetAll(ctx context.Context, t topic.ListParams) (pagination.List[topic.Topic], error)
}

type answerService interface {
	Create(ctx context.Context, a answer.Answer) (answer.Answer, error)
	GetAll(ctx context.Context, p answer.ListParams) (pagination.List[answer.Answer], error)
}

type questionService interface {
	Create(ctx context.Context, dto question.CreateDTO) (questionId uint32, err error)
	GetById(ctx context.Context, questionId uint32) (question.Detailed, error)
	GetAll(ctx context.Context, p question.ListParams) (pagination.List[question.Minimized], error)
}
