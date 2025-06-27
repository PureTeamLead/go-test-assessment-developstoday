package database

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfig struct {
	Database string `yaml:"database" env:"POSTGRES_DB"`
	Username string `yaml:"username" env:"POSTGRES_USER"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	Host     string `yaml:"host" env:"POSTGRES_HOST"`
	Port     int    `yaml:"port" env:"POSTGRES_PORT"`
	MinConns int    `yaml:"min-conns" env:"POSTGRES_MIN_CONNS"`
	MaxConns int    `yaml:"max-conns" env:"POSTGRES_MAX_CONNS"`
	MigPath  string `yaml:"mig-path" env:"POSTGRES_MIG_PATH"`
}

func NewPostgres(ctx context.Context, cfg PostgresConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	postgresCfg, err := pgxpool.ParseConfig(fmt.Sprintf("%s&pool_max_conns=%d&pool_min_conns=%d", connString, cfg.MaxConns, cfg.MinConns))
	if err != nil {
		return nil, fmt.Errorf("unable to parse database connection config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, postgresCfg)

	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err = doMigrate(connString, cfg.MigPath); err != nil {
		return nil, err
	}

	return pool, nil
}

func doMigrate(connStr, migPath string) error {
	m, err := migrate.New(fmt.Sprintf("file://%s", migPath), connStr)
	if err != nil {
		return fmt.Errorf("failed creating migrations: %w", err)
	}

	if err = m.Up(); err != nil {
		return fmt.Errorf("failed executing migrations: %w", err)
	}

	return nil
}
