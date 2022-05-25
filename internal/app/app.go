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

	"github.com/answersuck/vault/internal/adapter/repository"
	"github.com/answersuck/vault/internal/adapter/smtp"
	"github.com/answersuck/vault/internal/adapter/storage"
	v1 "github.com/answersuck/vault/internal/handler/http/v1"

	"github.com/answersuck/vault/internal/app/auth"
	"github.com/answersuck/vault/internal/app/email"

	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/internal/domain/language"
	"github.com/answersuck/vault/internal/domain/media"
	"github.com/answersuck/vault/internal/domain/question"
	"github.com/answersuck/vault/internal/domain/session"
	"github.com/answersuck/vault/internal/domain/tag"
	"github.com/answersuck/vault/internal/domain/topic"

	"github.com/answersuck/vault/pkg/blocklist"
	"github.com/answersuck/vault/pkg/httpserver"
	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/postgres"
	"github.com/answersuck/vault/pkg/token"
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
	sessionRepo := repository.NewSessionPSQL(l, pg)
	sessionService := session.NewService(&cfg.Session, sessionRepo)

	emailClient, err := smtp.NewClient(&smtp.ClientOptions{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		From:     cfg.SMTP.From,
		Password: cfg.SMTP.Password,
	})
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - email.NewClient: %w", err))
	}

	emailService := email.NewService(&cfg, emailClient)

	tokenManager, err := token.NewManager(cfg.AccessToken.Sign)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - auth.NewTokenManager: %w", err))
	}

	usernameBlockList := blocklist.New(blocklist.WithUsernames)

	accountRepo := repository.NewAccountPSQL(l, pg)
	accountService := account.NewService(&account.Deps{
		Config:         &cfg,
		AccountRepo:    accountRepo,
		SessionService: sessionService,
		EmailService:   emailService,
		Blocklist:      usernameBlockList,
	})

	authService := auth.NewService(&auth.Deps{
		Config:         &cfg,
		Token:          tokenManager,
		AccountService: accountService,
		SessionService: sessionService,
	})

	languageRepo := repository.NewLanguagePSQL(l, pg)
	languageService := language.NewService(languageRepo)

	tagRepo := repository.NewTagPSQL(l, pg)
	tagService := tag.NewService(tagRepo)

	topicRepo := repository.NewTopicPSQL(l, pg)
	topicService := topic.NewService(topicRepo)

	questionRepo := repository.NewQuestionPSQL(l, pg)
	questionService := question.NewService(questionRepo)

	ginTranslator, err := validation.NewGinTranslator()
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - validation.NewGinTranslator: %w", err))
	}

	storageProvider, err := storage.NewProvider(&cfg.FileStorage)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - storage.NewProvider: %w", err))
	}

	mediaRepo := repository.NewMediaPSQL(l, pg)
	mediaService := media.NewService(mediaRepo, storageProvider)

	// HTTP Server
	engine := gin.New()
	v1.NewHandler(
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
			MediaService:    mediaService,
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
