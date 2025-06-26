package server

import (
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Host         string        `yaml:"host" env:"HTTP_SERVER_HOST"`
	Port         int           `yaml:"port" env:"HTTP_SERVER_PORT"`
	ReadTimeout  time.Duration `yaml:"w-timeout" env:"HTTP_SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `yaml:"r-timeout" env:"HTTP_SERVER_WRITE_TIMEOUT"`
}

// TODO: maybe set a handler

func New(cfg *Config) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}
