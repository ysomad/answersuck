package psql

import (
	"log"
	"os"
	"testing"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/pkg/logger"
	"github.com/answersuck/vault/pkg/postgres"
)

const configPath = "../../../../configs/test.yml"

var (
	_testClient *postgres.Client
	_testCfg    *config.TestPG
	_testLogger *zap.Logger
)

func TestMain(m *testing.M) {
	var cfg config.Test

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Config error: %v", err)
	}
	_testCfg = &cfg.PG

	pg, err := postgres.NewClient(_testCfg.URL, postgres.MaxPoolSize(_testCfg.PoolMax))
	if err != nil {
		log.Fatalf("postgres.NewClient: %v", err)
	}
	_testClient = pg
	defer _testClient.Close()

	_testLogger = logger.New(os.Stdout, "")
	defer _testLogger.Sync()

	os.Exit(m.Run())
}
