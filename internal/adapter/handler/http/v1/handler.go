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
}

func NewMux(d *Deps) *chi.Mux {
	m := chi.NewMux()
	m.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	m.Mount("/accounts", newAccountHandler(d))
	m.Mount("/sessions", newSessionHandler(d))
	m.Mount("/auth", newAuthHandler(d))
	m.Mount("/media", newMediaHandler(d))

	return m
}

// func NewRouter(d *Deps) *fiber.App {
// 	r := fiber.New()

// 	r.Mount("/sessions", newSessionRouter(d))
// 	r.Mount("/accounts", newAccountRouter(d))
// 	r.Mount("/auth", newAuthRouter(d))
// 	r.Mount("/media", newMediaRouter(d))
// 	r.Mount("/languages", newLanguageRouter(d))
// 	r.Mount("/tags", newTagRouter(d))
// 	r.Mount("/topics", newTopicRouter(d))
// 	r.Mount("/answers", newAnswerRouter(d))
// 	r.Mount("/questions", newQuestionRouter(d))

// 	return r
// }
