package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/handler/http"
	v1 "github.com/answersuck/vault/internal/handler/http/v1"

	"github.com/answersuck/vault/internal/adapter/repository/psql"
	"github.com/answersuck/vault/internal/adapter/smtp"

	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/internal/domain/auth"
	"github.com/answersuck/vault/internal/domain/email"
	"github.com/answersuck/vault/internal/domain/session"

	"github.com/answersuck/vault/pkg/blocklist"
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
		l.Fatal(fmt.Errorf("main - run - postgres.NewClient: %w", err))
	}
	defer pg.Close()

	// Service
	validationModule, err := validation.NewModule()
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - validation.NewModule: %w", err))
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
		l.Fatal(fmt.Errorf("app - Run - email.NewClient: %w", err))
	}

	emailService := email.NewService(&cfg, emailClient)

	tokenManager, err := token.NewManager(cfg.AccessToken.Sign)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - token.NewManager: %w", err))
	}

	usernameBlockList := blocklist.New(blocklist.WithUsernames)

	accountRepo := psql.NewAccountRepo(l, pg)
	accountService := account.NewService(&account.Deps{
		Config:         &cfg,
		AccountRepo:    accountRepo,
		SessionService: sessionService,
		EmailService:   emailService,
		BlockList:      usernameBlockList,
	})

	authService := auth.NewService(&auth.Deps{
		Logger:         l,
		Config:         &cfg,
		Token:          tokenManager,
		AccountService: accountService,
		SessionService: sessionService,
	})

	// languageRepo := psql.NewLanguageRepo(l, pg)
	// languageService := language.NewService(languageRepo)

	// tagRepo := psql.NewTagRepo(l, pg)
	// tagService := tag.NewService(tagRepo)

	// topicRepo := psql.NewTopicRepo(l, pg)
	// topicService := topic.NewService(topicRepo)

	// questionRepo := psql.NewQuestionRepo(l, pg)
	// questionService := question.NewService(questionRepo)

	// storageProvider, err := storage.NewProvider(&cfg.FileStorage)
	// if err != nil {
	// 	l.Fatal(fmt.Errorf("app - Run - storage.NewProvider: %w", err))
	// }

	// mediaRepo := psql.NewMediaRepo(l, pg)
	// mediaService := media.NewService(mediaRepo, storageProvider)

	// answerRepo := psql.NewAnswerRepo(l, pg)
	// answerService := answer.NewService(l, answerRepo, mediaService)

	// playerRepo := psql.NewPlayerRepo(l, pg)
	// playerService := player.NewService(playerRepo)

	app := http.NewApp(cfg)

	app.Mount("/v1", v1.NewRouter(&v1.Deps{
		Config:           &cfg,
		Logger:           l,
		ValidationModule: validationModule,
		SessionService:   sessionService,
		AccountService:   accountService,
		AuthService:      authService,
	}))

	http.ServeSwaggerUI(app, cfg.HTTP.Debug)

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(cfg.HTTP.Port); err != nil {
			l.Fatal("app.Listen: %w", err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-c // This blocks the main thread until an interrupt is received
	l.Info("Gracefully shutting down...")

	_ = app.Shutdown()

	l.Info("HTTP App was successful shutdown.")
}
