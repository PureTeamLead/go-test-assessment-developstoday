package config

import (
	"fmt"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/transport/server"
	"github.com/PureTeamLead/go-test-assessment-developstoday/pkg/storage/postgres"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
)

const (
	envFilename = ".env"
)

type AppConfig struct {
	Env           string                  `yaml:"env" env:"ENVIRONMENT"`
	DBConfig      database.PostgresConfig `yaml:"db"`
	HTTPSrvConfig server.Config           `yaml:"http-server"`
}

func Load(path string) (*AppConfig, error) {
	const op = "config.Load"
	var cfg AppConfig

	err := godotenv.Load(envFilename)
	if err != nil {
		log.Println("[WARNING] .env file not found")
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("%s: failed to fill config data: %w", op, err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("%s: failed to read envs: %w", op, err)
	}

	return &cfg, nil
}
