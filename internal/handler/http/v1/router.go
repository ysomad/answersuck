package v1

import (
	"github.com/gofiber/fiber/v2"

	"github.com/answersuck/vault/internal/config"

	"github.com/answersuck/vault/pkg/logging"
)

type Deps struct {
	Config              *config.Aggregate
	Logger              logging.Logger
	ValidationModule    ValidationModule
	SessionService      SessionService
	AccountService      AccountService
	VerificationService VerificationService
	LoginService        LoginService
	TokenService        TokenService
	MediaService        MediaService
	LanguageService     LanguageService
	TagService          TagService
	TopicService        TopicService
}

type listResp struct {
	Result any `json:"result"`
}

func NewRouter(d *Deps) *fiber.App {
	r := fiber.New()

	r.Mount("/sessions", newSessionRouter(d))
	r.Mount("/accounts", newAccountRouter(d))
	r.Mount("/auth", newAuthRouter(d))
	r.Mount("/media", newMediaRouter(d))
	r.Mount("/languages", newLanguageRouter(d))
	r.Mount("/tags", newTagRouter(d))
	r.Mount("/topics", newTopicRouter(d))

	return r
}
