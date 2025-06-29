package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Bot struct {
		Token         string  `yaml:"token"`
		Whitelist     []int64 `yaml:"whitelist"`
		NotifyUserIDs []int64 `yaml:"notify_user_ids"` // Список ID пользователей для уведомлений
	} `yaml:"bot"`

	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
}

func New(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	cfg := new(Config)

	if err = yaml.Unmarshal(file, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	return cfg, nil
}
