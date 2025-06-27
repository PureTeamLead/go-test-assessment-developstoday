package config

import (
	"fmt"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/transport/server"
	"github.com/PureTeamLead/go-test-assessment-developstoday/pkg/storage/postgres"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
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

	_ = godotenv.Load(envFilename)
	//if err != nil {
	//	return nil, fmt.Errorf("%s: failed to load envs from env file: %w", op, err)
	//}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("%s: failed to fill config data: %w", op, err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("%s: failed to read env file: %w", op, err)
	}

	return &cfg, nil
}
