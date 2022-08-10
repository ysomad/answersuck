package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/answersuck/host/internal/config"
)

type Deps struct {
	Config          *config.Aggregate
	Logger          *zap.Logger
	Validate        validate
	AccountService  accountService
	SessionService  sessionService
	LoginService    loginService
	TokenService    tokenService
	MediaService    mediaService
	LanguageService languageService
	TagService      tagService
	TopicService    topicService
	AnswerService   answerService
	QuestionService questionService
	PlayerService   playerService
}

func NewMux(d *Deps) *chi.Mux {
	m := chi.NewMux()
	m.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	m.Mount("/accounts", newAccountMux(d))
	m.Mount("/sessions", newSessionMux(d))
	m.Mount("/auth", newAuthMux(d))
	m.Mount("/media", newMediaMux(d))
	m.Mount("/languages", newLanguageMux(d))
	m.Mount("/tags", newTagMux(d))
	m.Mount("/answers", newAnswerMux(d))
	m.Mount("/topics", newTopicMux(d))
	m.Mount("/questions", newQuestionMux(d))
	m.Mount("/players", newPlayerMux(d))

	return m
}
