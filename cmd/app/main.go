package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/Quizish/quizish-backend/internal/app"
	v1 "github.com/Quizish/quizish-backend/internal/handler/http/v1"
	"github.com/Quizish/quizish-backend/internal/repository"
	"github.com/Quizish/quizish-backend/internal/service"
	"github.com/Quizish/quizish-backend/pkg/auth"
	"github.com/Quizish/quizish-backend/pkg/httpserver"
	"github.com/Quizish/quizish-backend/pkg/logging"
	"github.com/Quizish/quizish-backend/pkg/mongodb"
	"github.com/Quizish/quizish-backend/pkg/postgres"
	"github.com/Quizish/quizish-backend/pkg/validation"
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

// run creates objects via constructors.
func run(cfg *app.Config) {
	l := logging.NewLogger(cfg.LogLevel)

	// Postgres
	pg, err := postgres.New(cfg.PostgresURL, postgres.MaxPoolSize(cfg.PostgresPoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("main - run - postgres.New: %w", err))
	}
	defer pg.Close()

	// MongoDB
	mcli, err := mongodb.NewClient(cfg.MongoURI, cfg.MongoUsername, cfg.MongoPassword)
	if err != nil {
		l.Fatal(fmt.Errorf("main - run - mongodb.NewClient: %w", err))
	}
	mdb := mcli.Database(cfg.MongoDatabase)

	/*
		rdb := redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       0,
		})
	*/

	// Service
	sessionRepo := repository.NewSessionRepository(mdb)
	sessionService := service.NewSessionService(cfg, sessionRepo)

	accountRepo := repository.NewAccountRepository(pg)
	accountService := service.NewAccountService(cfg, accountRepo, sessionService)

	tokenManager, err := auth.NewTokenManager(cfg.AccessTokenSigningKey)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - auth.NewJWTManager: %w", err))
	}

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
