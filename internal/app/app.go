package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/answersuck/vault/internal/config"

	v1 "github.com/answersuck/vault/internal/handler/http/v1"

	"github.com/answersuck/vault/internal/adapter/repository/psql"
	"github.com/answersuck/vault/internal/adapter/smtp"

	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/internal/domain/auth"
	"github.com/answersuck/vault/internal/domain/email"
	"github.com/answersuck/vault/internal/domain/session"

	"github.com/answersuck/vault/pkg/blocklist"
	"github.com/answersuck/vault/pkg/crypto"
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

	// DB
	pg, err := postgres.NewClient(cfg.PG.URL,
		postgres.MaxPoolSize(cfg.PG.PoolMax),
		postgres.PreferSimpleProtocol(cfg.PG.SimpleProtocol))
	if err != nil {
		l.Fatal("app - run - postgres.NewClient: %w", err)
	}
	defer pg.Close()

	// Service
	validationModule, err := validation.NewModule()
	if err != nil {
		l.Fatal("app - Run - validation.NewModule: %w", err)
	}

	sessionRepo := psql.NewSessionRepo(l, pg)
	sessionService := session.NewService(&cfg.Session, sessionRepo)

	emailClient, err := smtp.NewClient(&smtp.ClientOptions{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		From:     cfg.SMTP.From,
		Password: cfg.SMTP.Password,
	})
	if err != nil {
		l.Fatal("app - Run - email.NewClient: %w", err)
	}

	emailService := email.NewService(&cfg, emailClient)

	tokenManager, err := token.NewManager(cfg.SecurityToken.Sign)
	if err != nil {
		l.Fatal("app - Run - token.NewManager: %w", err)
	}

	usernameBlockList := blocklist.New(blocklist.WithUsernames)
	argon2Hasher := crypto.NewArgon2Id()

	accountRepo := psql.NewAccountRepo(l, pg)
	accountService := account.NewService(&account.Deps{
		Config:         &cfg,
		Password:       argon2Hasher,
		AccountRepo:    accountRepo,
		SessionService: sessionService,
		EmailService:   emailService,
		BlockList:      usernameBlockList,
	})

	loginService := auth.NewLoginService(accountService, sessionService, argon2Hasher)
	tokenService := auth.NewTokenService(auth.TokenServiceDeps{
		Config:          &cfg.SecurityToken,
		TokenManager:    tokenManager,
		AccountService:  accountService,
		PasswordMatcher: argon2Hasher,
	})
	//
	// languageRepo := psql.NewLanguageRepo(l, pg)
	// languageService := language.NewService(languageRepo)
	//
	// tagRepo := psql.NewTagRepo(l, pg)
	// tagService := tag.NewService(tagRepo)
	//
	// topicRepo := psql.NewTopicRepo(l, pg)
	// topicService := topic.NewService(topicRepo)
	//
	// questionRepo := psql.NewQuestionRepo(l, pg)
	// questionService := question.NewService(questionRepo)
	//
	// storageProvider, err := storage.NewProvider(&cfg.FileStorage)
	// if err != nil {
	// 	l.Fatal("app - Run - storage.NewProvider: %w", err)
	// }
	//
	// mediaRepo := psql.NewMediaRepo(l, pg)
	// mediaService := media.NewService(mediaRepo, storageProvider)
	//
	// answerRepo := psql.NewAnswerRepo(l, pg)
	// answerService := answer.NewService(l, answerRepo, mediaService)
	//
	// playerRepo := psql.NewPlayerRepo(l, pg)
	// playerService := player.NewService(playerRepo)

	// http
	m := chi.NewMux()

	m.Mount("/v1", v1.NewHandler(&v1.Deps{
		Config:           &cfg,
		Logger:           l,
		ValidationModule: validationModule,
		AccountService:   accountService,
		SessionService:   sessionService,
		LoginService:     loginService,
		TokenService:     tokenService,
	}))

	httpServer := httpserver.New(m, httpserver.Port(cfg.HTTP.Port))

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
