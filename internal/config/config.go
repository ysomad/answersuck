package config

import "time"

type Config struct {
	App     App     `yaml:"app"`
	HTTP    HTTP    `yaml:"http"`
	Log     Log     `yaml:"log"`
	PG      PG      `yaml:"postgres"`
	Session Session `yaml:"session"`
}

type App struct {
	Name string `yaml:"name" env-required:"true"`
	Ver  string `yaml:"version" env-required:"true"`
	Env  string `yaml:"environment" env-required:"true"`
}

type (
	HTTP struct {
		Host string `yaml:"host" env-required:"true"`
		Port string `yaml:"port" env-required:"true"`
	}

	Log struct {
		Level string `yaml:"level" env-required:"true"`
	}

	PG struct {
		MaxConns int32  `yaml:"max_connections" env-required:"true"`
		URL      string `env:"PG_URL" env-required:"true"`
		DBName   string `env:"PG_DB_NAME" env-required:"true"`
	}

	Session struct {
		LifeTime time.Duration `yaml:"lifetime" env-required:"true"`
	}
)
