package config

type (
	Test struct {
		PG TestPG `yaml:"postgres"`
	}

	TestPG struct {
		PoolMax int    `env-required:"true" yaml:"poolMax"`
		URL     string `env-required:"true" yaml:"url"`
	}
)
