package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/auth"

	"github.com/answersuck/vault/pkg/logging"
)

const route = "/api/v1"

type Deps struct {
	Config          *config.Aggregate
	Logger          logging.Logger
	ErrorTranslator ErrorTranslator
	TokenManager    auth.TokenManager
	AccountService  AccountService
	SessionService  SessionService
	AuthService     AuthService
	LanguageService LanguageService
	QuestionService QuestionService
	TagService      TagService
	TopicService    topicService
}

func NewHandler(r *gin.Engine, d *Deps) {
	// Options
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Resource handlers
	h := r.Group(route)
	{
		newSessionHandler(h, d)
		newAccountHandler(h, d)
		newAuthHandler(h, d)
		newLanguageHandler(h, d)
		newTagHandler(h, d)
		newTopicHandler(h, d)
		newQuestionHandler(h, d)
	}
}
