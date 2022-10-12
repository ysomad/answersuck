package config

type Config struct {
	App      App      `yaml:"app"`
	HTTP     HTTP     `yaml:"http"`
	Log      Log      `yaml:"log"`
	Postgres Postgres `yaml:"postgres"`
}

type (
	App struct {
		Name string `yaml:"name" env-required:"true"`
		Ver  string `yaml:"version" env-required:"true"`
	}

	HTTP struct {
		Port string `yaml:"port" env-required:"true"`
	}

	Log struct {
		Level string `yaml:"level" env-required:"true"`
	}

	Postgres struct {
		MaxConns int    `yaml:"max_connections" env-required:"true"`
		URL      string `env:"POSTGRES_URL" env-required:"true"`
	}
)
