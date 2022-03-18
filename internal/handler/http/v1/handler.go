package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quizlyfun/quizly-backend/internal/app"
	"github.com/quizlyfun/quizly-backend/internal/service"

	"github.com/quizlyfun/quizly-backend/pkg/auth"
	"github.com/quizlyfun/quizly-backend/pkg/logging"
	"github.com/quizlyfun/quizly-backend/pkg/validation"
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

// urlParam returns url parameter as a string for gin handler
func urlParam(param string) string {
	return fmt.Sprintf(":%s", param)
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
	}
}
