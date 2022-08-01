package config

import (
	"time"
)

type (
	Aggregate struct {
		App           App           `yaml:"app"`
		Web           Web           `yaml:"web"`
		HTTP          HTTP          `yaml:"http"`
		Log           Log           `yaml:"logger"`
		PG            PG            `yaml:"postgres"`
		Cache         Cache         `yaml:"cache"`
		SecurityToken SecurityToken `yaml:"securityToken"`
		Session       Session       `yaml:"session"`
		FileStorage   FileStorage   `yaml:"fileStorage"`
		Email         Email         `yaml:"email"`
		SMTP          SMTP          `yaml:"smtp"`
		Password      Password      `yaml:"password"`
	}

	Web struct {
		URL string `yaml:"url" env-required:"true"`
	}

	App struct {
		Name string `yaml:"name" env-required:"true"`
		Ver  string `yaml:"version" env-required:"true"`
	}

	HTTP struct {
		Port  string `yaml:"port" env-required:"true"`
		Debug bool   `yaml:"debug"`
	}

	Log struct {
		Level string `yaml:"logLevel" env-required:"true"`
	}

	PG struct {
		PoolMax        int    `yaml:"poolMax" env-required:"true"`
		URL            string `env:"PG_URL" env-required:"true"`
		SimpleProtocol bool   `yaml:"simpleProtocol"`
	}

	Cache struct {
		Expiration time.Duration `yaml:"expiration" env-required:"true"`
		DB         int           `yaml:"db"`
	}

	Session struct {
		Expiration     time.Duration `yaml:"expiration" env-required:"true"`
		CookieName     string        `yaml:"cookieName" env-required:"true"`
		CookieSecure   bool          `yaml:"cookieSecure"`
		CookieHTTPOnly bool          `yaml:"cookieHttpOnly"`
		CookiePath     string        `yaml:"cookiePath"`
	}

	SecurityToken struct {
		Expiration time.Duration `yaml:"expiration" env-required:"true"`
		Sign       string        `yaml:"signingKey" env:"ACCESS_TOKEN_SIGNING_KEY" env-required:"true"`
	}

	FileStorage struct {
		Endpoint  string `yaml:"endpoint" env-required:"true"`
		Bucket    string `yaml:"bucket" env-required:"true"`
		AccessKey string `env:"FILE_STORAGE_ACCESS_KEY" env-required:"true"`
		SecretKey string `env:"FILE_STORAGE_SECRET_KEY" env-required:"true"`
		Host      string `yaml:"host" env-required:"true"`
		CDNHost   string `yaml:"cdnHost" env-required:"true"`
		CDN       bool   `yaml:"cdn"`
		SSL       bool   `yaml:"ssl"`
	}

	Email struct {
		Template EmailTemplate `yaml:"templates" env-required:"true"`
		Subject  EmailSubject  `yaml:"subjects" env-required:"true"`
		Format   EmailFormat   `yaml:"formats" env-required:"true"`
	}

	EmailTemplate struct {
		AccountVerification  string `yaml:"accountVerification" env-required:"true"`
		AccountPasswordReset string `yaml:"accountPasswordReset" env-required:"true"`
	}

	EmailSubject struct {
		AccountVerification  string `yaml:"accountVerification" env-required:"true"`
		AccountPasswordReset string `yaml:"accountPasswordReset" env-required:"true"`
	}

	EmailFormat struct {
		AccountVerification  string `yaml:"accountVerification" env-required:"true"`
		AccountPasswordReset string `yaml:"accountPasswordReset" env-required:"true"`
	}

	SMTP struct {
		Host     string `yaml:"host" env:"SMTP_HOST" env-required:"true"`
		Port     int    `yaml:"port" env:"SMTP_PORT" env-required:"true"`
		From     string `env:"SMTP_FROM" env-required:"true"`
		Password string `env:"SMTP_PASSWORD" env-required:"true"`
	}

	Password struct {
		ResetTokenExpiration time.Duration `yaml:"resetTokenExpiration" env-required:"true"`
	}
)
