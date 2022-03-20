package app

import (
	"time"
)

type (
	Config struct {
		App         `yaml:"app"`
		HTTP        `yaml:"http"`
		Log         `yaml:"logger"`
		PG          `yaml:"postgres"`
		Cookie      `yaml:"cookie"`
		Cache       `yaml:"cache"`
		Redis       `yaml:"redis"`
		AccessToken `yaml:"access_token"`
		Session     `yaml:"session"`
		Email       `yaml:"email"`
		SMTP        `yaml:"smtp"`
	}

	App struct {
		AppName string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		AppVer  string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		HTTPPort string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		LogLevel string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
	}

	PG struct {
		PostgresPoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		PostgresURL     string `env-required:"true" env:"PG_URL"`
	}

	Cache struct {
		CacheTTL time.Duration `env-required:"true" yaml:"ttl" env:"CACHE_TTL"`
	}

	Redis struct {
		RedisAddr     string `env-required:"true" env:"REDIS_ADDR"`
		RedisPassword string `env-required:"true" env:"REDIS_PASSWORD"`
	}

	Cookie struct {
		CookieSecure   bool `yaml:"secure" env:"COOKIE_SECURE"`
		CookieHTTPOnly bool `yaml:"httponly" env:"COOKIE_HTTP_ONLY"`
	}

	Session struct {
		SessionTTL    time.Duration `env-required:"true" yaml:"ttl" env:"SESSION_TTL"`
		SessionCookie string        `env-required:"true" yaml:"cookie_key" env:"SESSION_COOKIE"`
		SessionDB     int           `yaml:"db" env:"SESSION_DB"`
	}

	AccessToken struct {
		AccessTokenTTL        time.Duration `env-required:"true" yaml:"ttl" env:"ACCESS_TOKEN_TTL"`
		AccessTokenSigningKey string        `env-required:"true" yaml:"signing_key" env:"ACCESS_TOKEN_SIGNING_KEY"`
	}

	Email struct {
		EmailVerificationTemplate string `env-required:"true" yaml:"verification_template" env:"EMAIL_VERIFICATION_TEMPLATE"`
		EmailVerificationSubject  string `env-required:"true" yaml:"verification_subject" env:"EMAIL_VERIFICATION_SUBJECT"`
		EmailVerificationLink     string `env-required:"true" yaml:"verification_link" env:"EMAIL_VERIFICATION_LINK"`
	}

	SMTP struct {
		SMTPHost string `env-required:"true" yaml:"host" env:"SMTP_HOST"`
		SMTPPort int    `env-required:"true" yaml:"port" env:"SMTP_PORT"`
		SMTPFrom string `env-required:"true" env:"SMTP_FROM"`
		SMTPPass string `env-required:"true" env:"SMTP_PASSWORD"`
	}
)
