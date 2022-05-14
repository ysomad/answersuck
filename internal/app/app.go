package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/answersuck/vault/internal/config"
	v1 "github.com/answersuck/vault/internal/handler/http/v1"
	repository "github.com/answersuck/vault/internal/repository/psql"
	"github.com/answersuck/vault/internal/service"

	"github.com/answersuck/vault/pkg/auth"
	"github.com/answersuck/vault/pkg/blocklist"
	"github.com/answersuck/vault/pkg/email"
	"github.com/answersuck/vault/pkg/httpserver"
	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
	"github.com/answersuck/vault/pkg/validation"
)

func Run(configPath string) {
	var cfg config.Aggregate

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	l := logging.NewLogger(cfg.Log.Level)

	l.Info(fmt.Sprintf("%+v\n", cfg))

	// DB
	pg, err := postgres.NewClient(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("main - run - postgres.NewClient: %w", err))
	}
	defer pg.Close()

	// Service
	sessionRepo := repository.NewSessionRepo(l, pg)
	sessionService := service.NewSessionService(&cfg.Session, sessionRepo)

	emailClient, err := email.NewClient(cfg.SMTP.From, cfg.SMTP.Password, cfg.SMTP.Host, cfg.SMTP.Port)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - email.NewClient: %w", err))
	}

	emailService := service.NewEmailService(&cfg, emailClient)

	tokenManager, err := auth.NewTokenManager(cfg.AccessToken.Sign)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - auth.NewTokenManager: %w", err))
	}

	usernameBlockList := blocklist.New(blocklist.WithUsernames)

	accountRepo := repository.NewAccountRepo(l, pg)
	accountService := service.NewAccountService(&cfg, accountRepo, sessionService, tokenManager,
		emailService, usernameBlockList)

	authService := service.NewAuthService(&cfg, tokenManager, accountService, sessionService)

	languageRepo := repository.NewLanguageRepo(l, pg)
	languageService := service.NewLanguageService(languageRepo)

	tagRepo := repository.NewTagRepo(l, pg)
	tagService := service.NewTagService(tagRepo)

	topicRepo := repository.NewTopicRepo(l, pg)
	topicService := service.NewTopicService(topicRepo)

	questionRepo := repository.NewQuestionRepo(l, pg)
	questionService := service.NewQuestionService(questionRepo)

	ginTranslator, err := validation.NewGinTranslator()
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - validation.NewGinTranslator: %w", err))
	}

	// HTTP Server
	engine := gin.New()
	v1.SetupHandlers(
		engine,
		&v1.Deps{
			Config:          &cfg,
			Logger:          l,
			ErrorTranslator: ginTranslator,
			TokenManager:    tokenManager,
			AccountService:  accountService,
			SessionService:  sessionService,
			AuthService:     authService,
			LanguageService: languageService,
			TagService:      tagService,
			TopicService:    topicService,
			QuestionService: questionService,
		},
	)

	// Swagger UI
	engine.Static("docs", "third_party/swaggerui")

	httpServer := httpserver.New(engine, httpserver.Port(cfg.HTTP.Port))

	// Graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
