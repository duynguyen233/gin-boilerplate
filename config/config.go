package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		GRPC `yaml:"grpc"`
		Log  `yaml:"logger"`
		PG   `yaml:"postgres"`
		Cors `yaml:"cors"`
	}

	// App -.
	App struct {
		Name    string `yaml:"name"    env:"APP_NAME"`
		Version string `yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `yaml:"port" env:"HTTP_PORT"`
	}

	// GRPC -.
	GRPC struct {
		Port int `yaml:"port" env:"GRPC_PORT"`
	}

	// Log -.
	Log struct {
		Level string `yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"false"                 env:"PG_URL"`
	}

	// Cors -.
	Cors struct {
		AllowedOrigins []string `yaml:"allowed_origins"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	_, er := os.ReadFile("./config/.env")
	if er != nil {
		panic(er)
	}

	err := cleanenv.ReadConfig("./config/.env", cfg)

	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
