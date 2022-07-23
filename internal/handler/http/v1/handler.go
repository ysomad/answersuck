package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/answersuck/vault/internal/config"

	"github.com/answersuck/vault/pkg/logging"
)

type Deps struct {
	Config           *config.Aggregate
	Logger           logging.Logger
	ValidationModule ValidationModule
	AccountService   AccountService
	SessionService   SessionService
	LoginService     LoginService
	TokenService     TokenService
	MediaService     MediaService
	LanguageService  LanguageService
	TagService       TagService
	TopicService     TopicService
	AnswerService    AnswerService
	QuestionService  QuestionService
}

func mountMiddlewares(r chi.Router) {
	for _, mw := range []func(http.Handler) http.Handler{
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	} {
		r.Use(mw)
	}
}

func NewHandler(d *Deps) http.Handler {
	r := chi.NewRouter()
	mountMiddlewares(r)

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
