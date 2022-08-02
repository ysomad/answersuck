package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/answersuck/host/internal/config"
)

type Deps struct {
	Config           *config.Aggregate
	Logger           *zap.Logger
	ValidationModule ValidationModule
	AccountService   accountService
	SessionService   sessionService
	LoginService     loginService
	TokenService     tokenService
	MediaService     mediaService
	LanguageService  languageService
	TagService       tagService
	TopicService     topicService
	AnswerService    answerService
	QuestionService  questionService
}

func NewHandler(d *Deps) http.Handler {
	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	r.Mount("/accounts", newAccountHandler(d))
	r.Mount("/sessions", newSessionHandler(d))
	r.Mount("/auth", newAuthHandler(d))

	return r
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
