package v1

import (
	"github.com/gofiber/fiber/v2"

	"github.com/answersuck/vault/internal/config"

	"github.com/answersuck/vault/pkg/logging"
)

type Deps struct {
	Config           *config.Aggregate
	Logger           logging.Logger
	ValidationModule ValidationModule
	SessionService   SessionService
	AccountService   AccountService
	LoginService     LoginService
	TokenService     TokenService
}

func NewRouter(d *Deps) *fiber.App {
	r := fiber.New()

	r.Mount("/sessions", newSessionRouter(d))
	r.Mount("/accounts", newAccountRouter(d))
	r.Mount("/auth", newAuthRouter(d))

	return r
}
