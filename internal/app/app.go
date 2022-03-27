package app

import (
	"fmt"
	repository2 "github.com/quizlyfun/quizly-backend/internal/repository"
	"github.com/quizlyfun/quizly-backend/pkg/blocklist"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/quizlyfun/quizly-backend/internal/config"
	v1 "github.com/quizlyfun/quizly-backend/internal/handler/http/v1"
	"github.com/quizlyfun/quizly-backend/internal/service"
	"github.com/quizlyfun/quizly-backend/pkg/auth"
	"github.com/quizlyfun/quizly-backend/pkg/email"
	"github.com/quizlyfun/quizly-backend/pkg/httpserver"
	"github.com/quizlyfun/quizly-backend/pkg/logging"
	"github.com/quizlyfun/quizly-backend/pkg/postgres"
	"github.com/quizlyfun/quizly-backend/pkg/storage"
	"github.com/quizlyfun/quizly-backend/pkg/validation"
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

	sessionRdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Session.DB,
	})

	// Service
	sessionRepo := repository2.NewSessionRepository(sessionRdb)
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

	minioClient, err := minio.New(cfg.FileStorage.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.FileStorage.AccessKey, cfg.FileStorage.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - minio.New: %w", err))
	}

	fileStorage := storage.NewFileStorage(minioClient, cfg.FileStorage.Bucket, cfg.FileStorage.Endpoint)
	usernameBlockList := blocklist.New(blocklist.WithUsernames)

	accountRepo := repository2.NewAccountRepository(pg)
	accountService := service.NewAccountService(&cfg, accountRepo, sessionService, tokenManager,
		emailService, fileStorage, usernameBlockList)

	authService := service.NewAuthService(&cfg, tokenManager, accountService, sessionService)

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
		},
	)

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
