package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/config"

	"github.com/answersuck/vault/pkg/auth"
	"github.com/answersuck/vault/pkg/logging"
)

const route = "/api/v1"

type Deps struct {
	Config          *config.Aggregate
	Logger          logging.Logger
	GinTranslator   errorTranslator
	TokenManager    auth.TokenManager // TODO: fix
	AccountService  accountService
	SessionService  sessionService
	AuthService     authService
	LanguageService languageService
	QuestionService questionService
	TagService      tagService
	TopicService    topicService
}

func SetupHandlers(e *gin.Engine, d *Deps) {
	// Options
	e.Use(gin.Logger())
	e.Use(gin.Recovery())

	e.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Resource handlers
	h := e.Group(route)
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
