package config

import (
	"time"
)

type (
	Aggregate struct {
		App         App         `yaml:"app"`
		Web         Web         `yaml:"web"`
		HTTP        HTTP        `yaml:"http"`
		Log         Log         `yaml:"logger"`
		PG          PG          `yaml:"postgres"`
		Cookie      Cookie      `yaml:"cookie"`
		Cache       Cache       `yaml:"cache"`
		Redis       Redis       `yaml:"redis"`
		AccessToken AccessToken `yaml:"accessToken"`
		Session     Session     `yaml:"session"`
		FileStorage FileStorage `yaml:"fileStorage"`
		Email       Email       `yaml:"email"`
		SMTP        SMTP        `yaml:"smtp"`
		Password    Password    `yaml:"password"`
	}

	Web struct {
		URL string `env-required:"true" yaml:"url"`
	}

	App struct {
		Name string `env-required:"true" yaml:"name"`
		Ver  string `env-required:"true" yaml:"version"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"logLevel"`
	}

	PG struct {
		PoolMax int    `env-required:"true" yaml:"poolMax"`
		URL     string `env-required:"true" env:"PG_URL"`
	}

	Cache struct {
		Expiration time.Duration `env-required:"true" yaml:"expiration"`
		DB         int           `yaml:"db"`
	}

	Redis struct {
		Addr     string `env-required:"true" env:"REDIS_ADDR"`
		Password string `env-required:"true" env:"REDIS_PASSWORD"`
	}

	Cookie struct {
		Secure   bool `yaml:"secure"`
		HTTPOnly bool `yaml:"httpOnly"`
	}

	Session struct {
		Expiration time.Duration `env-required:"true" yaml:"expiration"`
		CookieKey  string        `env-required:"true" yaml:"cookieKey"`

		// DB is number of database inside redis
		DB int `yaml:"db"`
	}

	AccessToken struct {
		Expiration time.Duration `env-required:"true" yaml:"expiration"`
		Sign       string        `env-required:"true" yaml:"signingKey" env:"ACCESS_TOKEN_SIGNING_KEY"`
	}

	FileStorage struct {
		Endpoint  string `env-required:"true" yaml:"endpoint" env:"FILE_STORAGE_ENDPOINT"`
		Bucket    string `env-required:"true" yaml:"bucket" env:"FILE_STORAGE_BUCKET"`
		AccessKey string `env-required:"true" env:"FILE_STORAGE_ACCESS_KEY"`
		SecretKey string `env-required:"true" env:"FILE_STORAGE_SECRET_KEY"`
	}

	Email struct {
		Template EmailTemplate `env-required:"true" yaml:"templates"`
		Subject  EmailSubject  `env-required:"true" yaml:"subjects"`
	}

	EmailTemplate struct {
		AccountVerification  string `env-required:"true" yaml:"accountVerification"`
		AccountPasswordReset string `env-required:"true" yaml:"accountPasswordReset"`
	}

	EmailSubject struct {
		AccountVerification  string `env-required:"true" yaml:"accountVerification"`
		AccountPasswordReset string `env-required:"true" yaml:"accountPasswordReset"`
	}

	SMTP struct {
		Host     string `env-required:"true" yaml:"host" env:"SMTP_HOST"`
		Port     int    `env-required:"true" yaml:"port" env:"SMTP_PORT"`
		From     string `env-required:"true" env:"SMTP_FROM"`
		Password string `env-required:"true" env:"SMTP_PASSWORD"`
	}

	Password struct {
		ResetTokenExp time.Duration `env-required:"true" yaml:"resetTokenExpiration"`
	}
)
