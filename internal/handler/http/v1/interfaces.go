package v1

import (
	"context"

	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/internal/domain/answer"
	"github.com/answersuck/vault/internal/domain/language"
	"github.com/answersuck/vault/internal/domain/media"
	"github.com/answersuck/vault/internal/domain/question"
	"github.com/answersuck/vault/internal/domain/session"
	"github.com/answersuck/vault/internal/domain/tag"
	"github.com/answersuck/vault/internal/domain/topic"
)

type AccountService interface {
	Create(ctx context.Context, r account.CreateReq) (account.Account, error)
	Delete(ctx context.Context, accountId string) error

	RequestVerification(ctx context.Context, accountId string) error
	Verify(ctx context.Context, code string) error

	ResetPassword(ctx context.Context, login string) error
	SetPassword(ctx context.Context, token, password string) error
}

type SessionService interface {
	GetByIdWithDetails(ctx context.Context, sessionId string) (*session.WithAccountDetails, error)
	GetById(ctx context.Context, sessionId string) (*session.Session, error)
	GetAll(ctx context.Context, accountId string) ([]*session.Session, error)
	Terminate(ctx context.Context, sessionId string) error
	TerminateWithExcept(ctx context.Context, accountId, sessionId string) error
}

type (
	LoginService interface {
		Login(ctx context.Context, login, password string, d session.Device) (*session.Session, error)
	}

	TokenService interface {
		Create(ctx context.Context, accountId, password string) (string, error)
		Parse(ctx context.Context, token string) (string, error)
	}
)

type MediaService interface {
	UploadAndSave(ctx context.Context, dto *media.UploadDTO) (media.Media, error)
}

type LanguageService interface {
	GetAll(ctx context.Context) ([]*language.Language, error)
}

type TagService interface {
	CreateMultiple(ctx context.Context, r []tag.CreateReq) ([]*tag.Tag, error)
	GetAll(ctx context.Context) ([]*tag.Tag, error)
}

type TopicService interface {
	Create(ctx context.Context, req topic.CreateReq) (topic.Topic, error)
	GetAll(ctx context.Context) ([]*topic.Topic, error)
}

type AnswerService interface {
	Create(ctx context.Context, r answer.CreateReq) (answer.Answer, error)
}

type QuestionService interface {
	Create(ctx context.Context, q *question.Question) (*question.Question, error)
	GetById(ctx context.Context, questionId int) (*question.Detailed, error)
	GetAll(ctx context.Context) ([]question.Minimized, error)
}
