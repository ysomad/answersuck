package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/quizlyfun/quizly-backend/internal/app"
	v1 "github.com/quizlyfun/quizly-backend/internal/handler/http/v1"
	"github.com/quizlyfun/quizly-backend/internal/service"
	"github.com/quizlyfun/quizly-backend/internal/service/repository"

	"github.com/quizlyfun/quizly-backend/pkg/auth"
	"github.com/quizlyfun/quizly-backend/pkg/email"
	"github.com/quizlyfun/quizly-backend/pkg/httpserver"
	"github.com/quizlyfun/quizly-backend/pkg/logging"
	"github.com/quizlyfun/quizly-backend/pkg/postgres"
	"github.com/quizlyfun/quizly-backend/pkg/validation"
)

func main() {
	var cfg app.Config

	// TODO: read configuration from flag
	err := cleanenv.ReadConfig("./configs/local.yml", &cfg)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	run(&cfg)
}

// Run creates objects via constructors.
func run(cfg *app.Config) {
	l := logging.NewLogger(cfg.LogLevel)

	// Postgres
	pg, err := postgres.New(cfg.PostgresURL, postgres.MaxPoolSize(cfg.PostgresPoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("main - run - postgres.New: %w", err))
	}
	defer pg.Close()

	sessionRdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.SessionDB,
	})

	// Service
	sessionRepo := repository.NewSessionRepository(sessionRdb)
	sessionService := service.NewSessionService(cfg, sessionRepo)

	emailClient, err := email.NewClient(cfg.SMTPFrom, cfg.SMTPPass, cfg.SMTPHost, cfg.SMTPPort)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - email.NewClient: %w", err))
	}

	emailService := service.NewEmailService(cfg, emailClient)

	tokenManager, err := auth.NewTokenManager(cfg.AccessTokenSigningKey)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - auth.NewTokenManager: %w", err))
	}

	accountRepo := repository.NewAccountRepository(pg)
	accountService := service.NewAccountService(cfg, accountRepo, sessionService, tokenManager, emailService)

	authService := service.NewAuthService(cfg, tokenManager, accountService, sessionService)

	ginTranslator, err := validation.NewGinTranslator()
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - validation.NewGinTranslator: %w", err))
	}

	// HTTP Server
	handler := gin.New()
	v1.SetupHandlers(
		handler,
		&v1.Deps{
			Config:          cfg,
			Logger:          l,
			ErrorTranslator: ginTranslator,
			TokenManager:    tokenManager,
			AccountService:  accountService,
			SessionService:  sessionService,
			AuthService:     authService,
		},
	)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTPPort))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
