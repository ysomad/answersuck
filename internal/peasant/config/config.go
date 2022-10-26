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
		VerifCodeLifetime time.Duration `yaml:"verification_code_lifetime" env-default:"24h"`
	}

	Password struct {
		TokenLifetime time.Duration `yaml:"token_lifetime" env-default:"10m"`
	}
)
