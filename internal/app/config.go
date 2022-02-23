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
		MongoDB     `yaml:"mongodb"`
		Cookie      `yaml:"cookie"`
		Cache       `yaml:"cache"`
		Redis       `yaml:"redis"`
		AccessToken `yaml:"access_token"`
		Session     `yaml:"session"`
		OAuth       `yaml:"oauth"`
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

	MongoDB struct {
		MongoURI      string `env-required:"true" env:"MONGO_URI"`
		MongoUsername string `env-required:"true" env:"MONGO_USER"`
		MongoPassword string `env-required:"true" env:"MONGO_PASS"`
		MongoDatabase string `env-required:"true" yaml:"database" env:"MONGO_DATABASE"`
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

	OAuth struct {
		GitHubClientID     string `yaml:"github_client_id" env-required:"true" env:"GH_CLIENT_ID"`
		GitHubClientSecret string `env-required:"true" env:"GH_CLIENT_SECRET"`
		GitHubScope        string `yaml:"github_scope" env-required:"true" env:"GH_SCOPE"`

		GoogleClientID     string `yaml:"google_client_id" env-required:"true" env:"GOOGLE_CLIENT_ID"`
		GoogleClientSecret string `env-required:"true" env:"GOOGLE_CLIENT_SECRET"`
		GoogleScope        string `yaml:"google_scope" env-required:"true" env:"GOOGLE_SCOPE"`

		DiscordClientID     string `yaml:"discord_client_id" env-required:"true" env:"DISCORD_CLIENT_ID"`
		DiscordClientSecret string `env-required:"true" env:"DISCORD_CLIENT_SECRET"`
		DiscordScope        string `yaml:"discord_scope" env-required:"true" env:"DISCORD_SCOPE"`
	}

	Session struct {
		SessionTTL    time.Duration `env-required:"true" yaml:"ttl" env:"SESSION_TTL"`
		SessionCookie string        `env-required:"true" yaml:"cookie_key" env:"SESSION_COOKIE"`
	}

	AccessToken struct {
		AccessTokenTTL        time.Duration `env-required:"true" yaml:"ttl" env:"ACCESS_TOKEN_TTL"`
		AccessTokenSigningKey string        `env-required:"true" yaml:"signing_key" env:"ACCESS_TOKEN_SIGNING_KEY"`
	}
)
