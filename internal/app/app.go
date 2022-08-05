package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"

	v1 "github.com/answersuck/host/internal/adapter/handler/http/v1"
	"github.com/answersuck/host/internal/adapter/repository/psql"
	"github.com/answersuck/host/internal/adapter/smtp"
	"github.com/answersuck/host/internal/adapter/storage"
	"github.com/answersuck/host/internal/config"
	"github.com/answersuck/host/internal/domain/account"
	"github.com/answersuck/host/internal/domain/auth"
	"github.com/answersuck/host/internal/domain/email"
	"github.com/answersuck/host/internal/domain/media"
	"github.com/answersuck/host/internal/domain/session"
	"github.com/answersuck/host/internal/pkg/blocklist"
	"github.com/answersuck/host/internal/pkg/crypto"
	"github.com/answersuck/host/internal/pkg/httpserver"
	"github.com/answersuck/host/internal/pkg/logger"
	"github.com/answersuck/host/internal/pkg/migrate"
	"github.com/answersuck/host/internal/pkg/postgres"
	"github.com/answersuck/host/internal/pkg/token"
	"github.com/answersuck/host/internal/pkg/validate"
)

func init() { migrate.Up("migrations") }

func Run(configPath string) {
	var cfg config.Aggregate

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	l := logger.New(os.Stdout, cfg.Log.Level)
	defer l.Sync()

	// DB
	pg, err := postgres.NewClient(cfg.PG.URL,
		postgres.MaxPoolSize(cfg.PG.PoolMax),
		postgres.PreferSimpleProtocol(cfg.PG.SimpleProtocol))
	if err != nil {
		l.Fatal("app - run - postgres.NewClient", zap.Error(err))
	}
	defer pg.Close()

	// Service
	validate, err := validate.New()
	if err != nil {
		l.Fatal("app - Run - validate.New", zap.Error(err))
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
		l.Fatal("app - Run - email.NewClient", zap.Error(err))
	}

	emailService := email.NewService(&cfg, emailClient)

	tokenManager, err := token.NewManager(cfg.SecurityToken.Sign)
	if err != nil {
		l.Fatal("app - Run - token.NewManager", zap.Error(err))
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

	storageProvider, err := storage.NewProvider(&cfg.FileStorage)
	if err != nil {
		l.Fatal("app - Run - storage.NewProvider", zap.Error(err))
	}

	mediaRepo := psql.NewMediaRepo(l, pg)
	mediaService := media.NewService(mediaRepo, storageProvider)

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

	//
	// answerRepo := psql.NewAnswerRepo(l, pg)
	// answerService := answer.NewService(l, answerRepo, mediaService)
	//
	// playerRepo := psql.NewPlayerRepo(l, pg)
	// playerService := player.NewService(playerRepo)

	// http
	m := chi.NewMux()

	m.Mount("/v1", v1.NewMux(&v1.Deps{
		Config:         &cfg,
		Logger:         l,
		Validate:       validate,
		AccountService: accountService,
		SessionService: sessionService,
		LoginService:   loginService,
		TokenService:   tokenService,
		MediaService:   mediaService,
	}))

	httpServer := httpserver.New(m, httpserver.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error("app - Run - httpServer.Notify", zap.Error(err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error("app - Run - httpServer.Shutdown", zap.Error(err))
	}
}
