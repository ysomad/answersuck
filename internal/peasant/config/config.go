package config

import "time"

type Config struct {
	App      App      `yaml:"app"`
	HTTP     HTTP     `yaml:"http"`
	Log      Log      `yaml:"log"`
	PG       PG       `yaml:"postgres"`
	Email    Email    `yaml:"email"`
	Password Password `yaml:"password"`
}

type App struct {
	Name string `yaml:"name" env-required:"true"`
	Ver  string `yaml:"version" env-required:"true"`
}

// Issuer using in jwt tokens
func (a App) Issuer() string {
	return a.Name + "-" + a.Ver
}

type (
	HTTP struct {
		Port string `yaml:"port" env-required:"true"`
	}

	Log struct {
		Level   string         `yaml:"level" env-required:"true"`
		TimeLoc *time.Location `yaml:"time_location" env-default:"Etc/UTC"`
	}

	PG struct {
		MaxConns int32  `yaml:"max_connections" env-required:"true"`
		URL      string `env:"PG_URL" env-required:"true"`
	}

	Email struct {
		VerifTokenSecret string        `env:"EMAIL_VERIF_TOKEN_SECRET" env-required:"true"`
		VerifTokenExp    time.Duration `yaml:"verif_token_expiration" env-required:"true"`
	}

	Password struct {
		SetterTokenSecret string        `env:"SETTER_PASSWORD_TOKEN_SECRET" env-required:"true"`
		SetterTokenExp    time.Duration `yaml:"setter_token_expiration" env-required:"true"`
	}
)
