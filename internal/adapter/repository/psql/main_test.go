package psql_test

import (
	"log"
	"os"
	"testing"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/pkg/postgres"
)

const configPath = "./configs/local.yml"

var (
	testClient *postgres.Client
)

func TestMain(m *testing.M) {
	var cfg config.PG

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	pg, err := postgres.NewClient(cfg.URL,
		postgres.MaxPoolSize(cfg.PoolMax),
		postgres.PreferSimpleProtocol(cfg.SimpleProtocol))
	if err != nil {
		log.Fatalf("postgres.NewClient: %v", err)
	}
	defer pg.Close()

	testClient = pg

	os.Exit(m.Run())
}
