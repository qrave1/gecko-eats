package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Bot      Bot
	Database Database
}

type Bot struct {
	Debug       bool    `env:"BOT_DEBUG" envDefault:"false"`
	Token       string  `env:"BOT_TOKEN,required"`
	Whitelist   []int64 `env:"BOT_WHITELIST"`
	NotifyUsers []int64 `env:"BOT_NOTIFY_USERS"` // Список ID пользователей для уведомлений
}

type Database struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
	SSL      string `env:"DB_SSL"`
}

func New() (*Config, error) {
	cfg, err := env.ParseAs[Config]()

	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
