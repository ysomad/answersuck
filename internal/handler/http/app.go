package http

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/answersuck/vault/internal/config"
)

const (
	logFmt           = "${pid} ${locals:requestid} ${status} - ${method} ${path}\n"
	defaultBodyLimit = 5 << 20 // 5MB
)

func ServeSwaggerUI(app *fiber.App, debug bool) {
	if debug {
		app.Static("/docs", "./third_party/swaggerui")
	}
}

func NewApp(cfg *config.Aggregate) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:     cfg.App.Name,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		BodyLimit:   defaultBodyLimit,
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.HTTP.Debug,
	}))

	app.Use(logger.New(logger.Config{Format: logFmt}))

	app.Use(requestid.New())

	return app
}
