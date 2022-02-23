package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Quizish/quizish-backend/internal/app"
	"github.com/Quizish/quizish-backend/internal/service"
	"github.com/Quizish/quizish-backend/pkg/auth"
	"github.com/Quizish/quizish-backend/pkg/logging"
	"github.com/Quizish/quizish-backend/pkg/validation"
)

const route = "/api/v1"

type Deps struct {
	Config          *app.Config
	Logger          logging.Logger
	ErrorTranslator validation.ErrorTranslator
	TokenManager    auth.TokenManager
	AccountService  service.Account
	SessionService  service.Session
	AuthService     service.Auth
}

func SetupHandlers(handler *gin.Engine, d *Deps) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger UI
	handler.Static(fmt.Sprintf("%s/swagger/", route), "third_party/swaggerui")

	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Resource handlers
	h := handler.Group(route)
	{
		newAccountHandler(h, d)
		newAuthHandler(h, d)
		newSessionHandler(h, d)
	}
}
