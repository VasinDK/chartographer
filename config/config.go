package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App         `yaml:"app"`
	HTTP        `yaml:"http"`
	Logger      `yaml:"logger"`
	FileStorage `yaml:"filestorage"`
}

type App struct {
	Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
	Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
}

type HTTP struct {
	Port         string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	WaitingClose int    `yaml:"waiting_timeout_close" env:"WAIT_CLOSE"`
}

type Logger struct {
	Level string `env-required:"true" yaml:"log_env" env:"LOG_ENV"`
}

type FileStorage struct {
	Path string `env-required:"true" yaml:"path" env:"PATH_FILE_ENV"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yaml", cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("config error %w", err)
	}

	return cfg, nil
}
